package dirmerger

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"chast.io/core/internal/internal_util/collection"
	chastlog "chast.io/core/internal/logger"
	"github.com/joomcode/errorx"
	"github.com/spf13/afero"
)

var errMergeOverwriteBlock = errorx.InternalError.New(
	"Error due to attempting to merge a file over an existing file in blockOverwrite mode",
)

func AreMergeable(mergeEntities []MergeEntity, targetFolder string, options *MergeOptions) (bool, error) {
	augmentedMergeOptions := *options
	augmentedMergeOptions.DryRun = true
	augmentedMergeOptions.DeleteEmptyFolders = false
	augmentedMergeOptions.DeleteMarkedAsDeletedPaths = false

	mergeError := MergeFolders(mergeEntities, targetFolder, &augmentedMergeOptions)

	if mergeError != nil {
		if errors.Is(mergeError, errMergeOverwriteBlock) {
			return false, nil
		}

		return false, mergeError
	}

	return true, nil
}

func MergeFolders(mergeEntities []MergeEntity, targetFolder string, options *MergeOptions) error {
	if !options.DryRun {
		if err := os.MkdirAll(targetFolder, options.FolderPermission); err != nil {
			return errorx.ExternalError.Wrap(err, fmt.Sprintf("failed to create target folder \"%s\"", targetFolder))
		}
	}

	for _, mergeEntity := range mergeEntities {
		if err := mergeFolders(mergeEntity, targetFolder, options); err != nil {
			return errorx.InternalError.Wrap(err,
				fmt.Sprintf("failed to merge folder \"%s\" with \"%s\"", mergeEntity.SourcePath, targetFolder),
			)
		}
	}

	return nil
}

func mergeFolders(mergeEntity MergeEntity, targetFolder string, options *MergeOptions) error {
	entityMergeOptions := *options
	entityMergeOptions.Inclusions = append(
		entityMergeOptions.Inclusions,
		collection.Map(mergeEntity.ChangeLocations.Include, NewWildcardString)...,
	)
	entityMergeOptions.Exclusions = append(
		entityMergeOptions.Exclusions,
		collection.Map(mergeEntity.ChangeLocations.Exclude, NewWildcardString)...,
	)

	sourceFolder := mergeEntity.SourcePath

	if err := mergeSourceIntoTarget(sourceFolder, targetFolder, &entityMergeOptions); err != nil {
		return errorx.InternalError.Wrap(err, "failed to merge source into target")
	}

	if entityMergeOptions.MergeMetaFilesFolder {
		if err := flattenMetaFolder(targetFolder, &entityMergeOptions); err != nil {
			return errorx.InternalError.Wrap(err, "failed to merge meta folder into source")
		}
	}

	if entityMergeOptions.DeleteEmptyFolders {
		if err := removeEmptyFolders(targetFolder, &entityMergeOptions); err != nil {
			return errorx.InternalError.Wrap(err, "failed to remove empty folders")
		}
	}

	if entityMergeOptions.DeleteMarkedAsDeletedPaths {
		if err := removeMarkedAsDeletedPaths(targetFolder, &entityMergeOptions); err != nil {
			return errorx.InternalError.Wrap(err, "failed to remove marked as deleted paths")
		}
	}

	return nil
}

func mergeSourceIntoTarget(sourceFolder string, targetFolder string, options *MergeOptions) error {
	osFileSystem := afero.NewOsFs()

	if walkError := afero.Walk(osFileSystem, sourceFolder, func(path string, info fs.FileInfo, _ error) error {
		if info == nil {
			return nil // exit if the file does not exist - this can happen due to a removal in a previous iteration
		}

		if path == sourceFolder {
			return nil // skip the root folder
		}

		if options.ShouldSkip(strings.TrimPrefix(path, sourceFolder)) {
			chastlog.Log.Tracef("Skipping path \"%s\" due to being excluded or not included", path)

			return nil // exit if the path is to be skipped
		}

		if info.IsDir() {
			if err := moveFolder(path, sourceFolder, targetFolder, osFileSystem, options); err != nil {
				return errorx.InternalError.Wrap(err, "Failed to move folder")
			}
		} else {
			if err := moveFile(path, sourceFolder, targetFolder, osFileSystem, options); err != nil {
				return errorx.InternalError.Wrap(err, "Failed to move file")
			}
		}

		return nil
	}); walkError != nil {
		return errorx.ExternalError.Wrap(walkError, "Failed to walk through source folder")
	}

	if !options.DryRun && sourceFolder != targetFolder && !options.CopyMode {
		if err := cleanupPath(sourceFolder, osFileSystem, options); err != nil {
			return err
		}
	}

	return nil
}

func moveFolder(
	sourcePath string,
	sourceRootFolder string,
	targetRootFolder string,
	osFileSystem afero.Fs,
	options *MergeOptions,
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

	targetPath := targetPath(sourcePath, sourceRootFolder, targetRootFolder)

	if !strings.HasPrefix(strings.TrimPrefix(sourcePath, sourceRootFolder), "/"+options.MetaFilesLocation) { // TODO cleanup
		if err := handleConflictingFolder(sourcePath, targetPath, osFileSystem, options); err != nil {
			return err
		}
	}

	if !options.DryRun {
		if err := osFileSystem.MkdirAll(targetPath, options.FolderPermission); err != nil {
			return errorx.ExternalError.Wrap(err, fmt.Sprintf("Failed to create folder \"%s\"", targetPath))
		}

		if !options.CopyMode {
			if err := cleanupPath(sourcePath, osFileSystem, options); err != nil {
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
	options *MergeOptions,
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

	if options.DryRun {
		return nil
	}

	// cases:
	// 1. source folder does not exist -> copy
	// 2. folder -> deleted:
	//    a. target folder is empty, delete it
	//    b. target folder is not empty, rename it
	// 3. deleted -> folder:
	//    a. target folder is empty, rename it
	//    b. target folder is not empty, delete folder
	if conflictingPathIsEmpty || isDeletedPath {
		if err := cleanupPath(conflictingPath, osFileSystem, options); err != nil {
			return errorx.InternalError.Wrap(err, "Failed to cleanup conflicting path")
		}
	} else {
		if err := os.Rename(conflictingPath, targetPath); err != nil {
			return errorx.ExternalError.Wrap(err, fmt.Sprintf("Failed to rename conflicting path \"%s\"", conflictingPath))
		}
	}

	return nil
}

func moveFile(
	sourcePath string,
	sourceRootFolder string,
	targetRootFolder string,
	osFileSystem afero.Fs,
	options *MergeOptions,
) error {
	if exists, err := afero.Exists(osFileSystem, sourcePath); err != nil || !exists {
		return nil //nolint:nilerr // If the folder does not exist, ignore it and continue
	}

	targetPath := targetPath(sourcePath, sourceRootFolder, targetRootFolder)

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
	options *MergeOptions,
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
			return errorx.InternalError.Wrap(errMergeOverwriteBlock, "Failed to move path %s to %s", sourcePath, targetPath)
		}

		if !options.DryRun {
			if err := osFileSystem.RemoveAll(conflictingPath); err != nil {
				return errorx.ExternalError.Wrap(err, "Failed to remove original path")
			}
		}
	}

	return nil
}

func targetPath(path string, sourceFolder string, targetFolder string) string {
	correctedPath := strings.TrimPrefix(path, sourceFolder)
	targetPath := filepath.Join(targetFolder, correctedPath)

	return targetPath
}

func cleanupPath(path string, _ afero.Fs, _ *MergeOptions) error {
	if err := os.RemoveAll(path); err != nil {
		return errorx.ExternalError.Wrap(err, "failed to remove merge source directory")
	}

	return nil
}
