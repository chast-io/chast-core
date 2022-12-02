package foldermover

import (
	"fmt"
	"os"
	"strings"

	chastlog "chast.io/core/internal/logger"
	pathutils "chast.io/core/internal/post_processing/merger/internal/path_utils"
	"chast.io/core/internal/post_processing/merger/pkg/mergeoptions"
	"chast.io/core/internal/post_processing/merger/pkg/mergererrors"
	"github.com/joomcode/errorx"
	"github.com/spf13/afero"
)

func MoveFolder(
	sourcePath string,
	sourceRootFolder string,
	targetRootFolder string,
	osFileSystem afero.Fs,
	options *mergeoptions.MergeOptions,
) error {
	if exists, err := afero.Exists(osFileSystem, sourcePath); err != nil || !exists {
		return nil //nolint:nilerr // If the folder does not exist, ignore it and continue
	}

	isEmpty, isEmptyCheckError := afero.IsEmpty(osFileSystem, sourcePath)
	if isEmptyCheckError != nil {
		return errorx.ExternalError.Wrap(isEmptyCheckError, "Failed to check if folder is empty")
	}

	if !isEmpty {
		chastlog.Log.Tracef("Folder \"%s\" is not empty, skipping -> will be handled later", sourcePath)

		return nil
	}

	targetPath := pathutils.TargetPath(sourcePath, sourceRootFolder, targetRootFolder)

	if !pathutils.IsInMetaFolder(sourcePath, sourceRootFolder, options) {
		if err := handleConflictingFolder(sourcePath, targetPath, osFileSystem, options); err != nil {
			return err
		}
	}

	if !options.DryRun {
		if err := osFileSystem.MkdirAll(targetPath, options.FolderPermission); err != nil {
			return errorx.ExternalError.Wrap(err, fmt.Sprintf("Failed to create folder \"%s\"", targetPath))
		}

		if !options.CopyMode {
			if err := pathutils.CleanupPath(sourcePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func handleConflictingFolder(
	sourcePath string,
	targetPath string,
	osFileSystem afero.Fs,
	options *mergeoptions.MergeOptions,
) error {
	isDeletedPath := strings.HasSuffix(sourcePath, options.MetaFilesDeletedExtension)

	var conflictingPath string
	if isDeletedPath { // is a deleted path
		conflictingPath = strings.TrimSuffix(targetPath, options.MetaFilesDeletedExtension)
	} else { // is a normal path
		conflictingPath = targetPath + options.MetaFilesDeletedExtension
	}

	conflictingPathExists, conflictingPathExistsError := afero.Exists(osFileSystem, conflictingPath)
	if conflictingPathExistsError != nil {
		return errorx.ExternalError.Wrap(conflictingPathExistsError, "Failed to check if conflicting path exists")
	}

	if !conflictingPathExists {
		return nil // nothing to do
	}

	conflictingPathIsEmpty, conflictingPathIsEmptyError := afero.IsEmpty(osFileSystem, conflictingPath)
	if conflictingPathIsEmptyError != nil {
		return errorx.ExternalError.Wrap(conflictingPathIsEmptyError, "Failed to check if conflicting path is empty")
	}

	// Case 1: Source folder does not exist -> do nothing
	if options.DryRun {
		return nil
	}

	if options.BlockOverwrite {
		return errorx.InternalError.Wrap(mergererrors.ErrMergeOverwriteBlock,
			"Cannot overwrite conflicting path \"%s\" with \"%s\"", conflictingPath, sourcePath)
	}

	// Case 2: existing folder -> deleted path:
	//    a. target folder is empty, delete it
	//    b. target folder is not empty, rename it
	// Case 3: deleted path -> new folder:
	//    a. target folder is empty, delete it
	if conflictingPathIsEmpty || isDeletedPath {
		if err := pathutils.CleanupPath(conflictingPath); err != nil {
			return errorx.InternalError.Wrap(err, "Failed to cleanup conflicting path")
		}
	} else {
		// Case 3: deleted path -> new folder:
		// 	  b: target folder is not empty, delete folder
		if err := os.Rename(conflictingPath, targetPath); err != nil {
			return errorx.ExternalError.Wrap(err, fmt.Sprintf("Failed to rename conflicting path \"%s\"", conflictingPath))
		}
	}

	return nil
}
