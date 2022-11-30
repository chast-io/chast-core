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

	if entityMergeOptions.MergeMetaFilesFolder {
		if err := mergeMetaFolderIntoSource(sourceFolder, &entityMergeOptions); err != nil {
			return errorx.InternalError.Wrap(err, "failed to merge meta folder into source")
		}
	}

	if err := mergeSourceIntoTarget(sourceFolder, targetFolder, &entityMergeOptions); err != nil {
		return errorx.InternalError.Wrap(err, "failed to merge source into target")
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

func mergeMetaFolderIntoSource(sourceFolder string, options *MergeOptions) error {
	metaFolderOptions := *options
	metaFolderOptions.BlockOverwrite = true

	return mergeSourceIntoTarget(
		filepath.Join(sourceFolder, metaFolderOptions.MetaFilesLocation),
		sourceFolder,
		&metaFolderOptions,
	)
}

func mergeSourceIntoTarget(sourceFolder string, targetFolder string, options *MergeOptions) error {
	osFileSystem := afero.NewOsFs()

	if walkError := afero.Walk(osFileSystem, sourceFolder, func(path string, info fs.FileInfo, _ error) error {
		if info == nil {
			return nil // exit if the file does not exist - this can happen due to a removal in a previous iteration
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

	if err := handleConflictingMovePaths(sourcePath, targetPath, osFileSystem, options); err != nil {
		return err
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

	if err := handleConflictingMovePaths(sourcePath, targetPath, osFileSystem, options); err != nil {
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

func handleConflictingMovePaths(
	sourcePath string,
	targetPath string,
	osFileSystem afero.Fs,
	options *MergeOptions,
) error {
	var conflictingPath string

	if strings.HasSuffix(sourcePath, options.MetaFilesDeletedExtension) {
		conflictingPath = strings.TrimSuffix(targetPath, options.MetaFilesDeletedExtension)
	} else {
		conflictingPath = targetPath + options.MetaFilesDeletedExtension
	}

	existingCounterpartExists, existingCounterpartExistenceCheckError := afero.Exists(osFileSystem, targetPath)
	if existingCounterpartExistenceCheckError != nil {
		return errorx.ExternalError.Wrap(
			existingCounterpartExistenceCheckError,
			"Failed to check if counterpart exists [case - existing file]",
		)
	}

	undeletedCounterpartExists, undeletedCounterpartExistenceCheckError := afero.Exists(osFileSystem, conflictingPath)
	if undeletedCounterpartExistenceCheckError != nil {
		return errorx.ExternalError.Wrap(
			undeletedCounterpartExistenceCheckError,
			"Failed to check if counterpart exists [case - deleted file]",
		)
	}

	if existingCounterpartExists || undeletedCounterpartExists {
		if options.BlockOverwrite {
			return errorx.WithPayload(errMergeOverwriteBlock, struct {
				sourcePath      string
				conflictingPath string
			}{sourcePath, conflictingPath})
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
