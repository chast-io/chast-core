package filemover

import (
	"os"
	"path/filepath"
	"strings"

	pathutils "chast.io/core/internal/post_processing/merger/internal/path_utils"
	"chast.io/core/internal/post_processing/merger/pkg/mergeoptions"
	"chast.io/core/internal/post_processing/merger/pkg/mergererrors"
	"github.com/joomcode/errorx"
	"github.com/spf13/afero"
)

func MoveFile(
	sourcePath string,
	sourceRootFolder string,
	targetRootFolder string,
	osFileSystem afero.Fs,
	options *mergeoptions.MergeOptions,
) error {
	if exists, err := afero.Exists(osFileSystem, sourcePath); err != nil || !exists {
		return nil //nolint:nilerr // If the folder does not exist, ignore it and continue
	}

	targetPath := pathutils.TargetPath(sourcePath, sourceRootFolder, targetRootFolder)

	if err := handleConflictingFile(sourcePath, targetPath, osFileSystem, options); err != nil {
		return err
	}

	if options.DryRun {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), options.FolderPermission); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to create target directory")
	}

	if options.CopyMode {
		if err := os.Link(sourcePath, targetPath); err != nil {
			return errorx.ExternalError.Wrap(err, "Failed to link file")
		}
	} else {
		if err := osFileSystem.Rename(sourcePath, targetPath); err != nil {
			return errorx.ExternalError.Wrap(err, "Failed to move file")
		}
	}

	return nil
}

func handleConflictingFile(
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

	existingCounterpartExists, existingCounterpartExistenceCheckError := afero.Exists(osFileSystem, targetPath)
	if existingCounterpartExistenceCheckError != nil {
		return errorx.ExternalError.Wrap(
			existingCounterpartExistenceCheckError,
			"Failed to check if counterpart exists [case - existing file]",
		)
	}

	deletionCounterpartExists, deletionCounterpartExistenceCheckError := afero.Exists(osFileSystem, conflictingPath)
	if deletionCounterpartExistenceCheckError != nil {
		return errorx.ExternalError.Wrap(
			deletionCounterpartExistenceCheckError,
			"Failed to check if counterpart exists [case - deleted file]",
		)
	}

	if (!isDeletedPath && existingCounterpartExists) || deletionCounterpartExists {
		if options.BlockOverwrite {
			return errorx.InternalError.Wrap(mergererrors.ErrMergeOverwriteBlock,
				"Failed to move path %s to %s", sourcePath, targetPath)
		}

		if !options.DryRun {
			if err := osFileSystem.RemoveAll(conflictingPath); err != nil {
				return errorx.ExternalError.Wrap(err, "Failed to remove original path")
			}
		}
	}

	return nil
}
