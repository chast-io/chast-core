package dirmerger

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"chast.io/core/internal/internal_util/collection"
	wildcardstring "chast.io/core/internal/internal_util/wildcard_string"
	chastlog "chast.io/core/internal/logger"
	"chast.io/core/internal/post_processing/merger/internal/dirmerger"
	filemover "chast.io/core/internal/post_processing/merger/internal/file_mover"
	foldermover "chast.io/core/internal/post_processing/merger/internal/folder_mover"
	metaflatterner "chast.io/core/internal/post_processing/merger/internal/meta_flatterner"
	pathutils "chast.io/core/internal/post_processing/merger/internal/path_utils"
	"chast.io/core/internal/post_processing/merger/pkg/mergeoptions"
	"chast.io/core/internal/post_processing/merger/pkg/mergererrors"
	"github.com/joomcode/errorx"
	"github.com/spf13/afero"
)

func AreMergeable(mergeEntities []MergeEntity, targetFolder string, options *mergeoptions.MergeOptions) (bool, error) {
	augmentedMergeOptions := *options
	augmentedMergeOptions.DryRun = true
	augmentedMergeOptions.DeleteEmptyFolders = false
	augmentedMergeOptions.DeleteMarkedAsDeletedPaths = false

	mergeError := MergeFolders(mergeEntities, targetFolder, &augmentedMergeOptions)

	if mergeError != nil {
		if errors.Is(mergeError, mergererrors.ErrMergeOverwriteBlock) {
			return false, nil
		}

		return false, mergeError
	}

	return true, nil
}

func MergeFolders(mergeEntities []MergeEntity, targetFolder string, options *mergeoptions.MergeOptions) error {
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

func mergeFolders(mergeEntity MergeEntity, targetFolder string, options *mergeoptions.MergeOptions) error {
	entityMergeOptions := *options
	entityMergeOptions.Inclusions = append(
		entityMergeOptions.Inclusions,
		collection.Map(mergeEntity.ChangeLocations.Include, wildcardstring.NewWildcardString)...,
	)
	entityMergeOptions.Exclusions = append(
		entityMergeOptions.Exclusions,
		collection.Map(mergeEntity.ChangeLocations.Exclude, wildcardstring.NewWildcardString)...,
	)

	sourceFolder := mergeEntity.SourcePath

	if err := mergeSourceIntoTarget(sourceFolder, targetFolder, &entityMergeOptions); err != nil {
		return errorx.InternalError.Wrap(err, "failed to merge source into target")
	}

	if entityMergeOptions.MergeMetaFilesFolder {
		if err := metaflatterner.FlattenMetaFolder(targetFolder, &entityMergeOptions); err != nil {
			return errorx.InternalError.Wrap(err, "failed to merge meta folder into source")
		}
	}

	if entityMergeOptions.DeleteEmptyFolders {
		if err := dirmerger.RemoveEmptyFolders(targetFolder, &entityMergeOptions); err != nil {
			return errorx.InternalError.Wrap(err, "failed to remove empty folders")
		}
	}

	if entityMergeOptions.DeleteMarkedAsDeletedPaths {
		if err := dirmerger.RemoveMarkedAsDeletedPaths(targetFolder, &entityMergeOptions); err != nil {
			return errorx.InternalError.Wrap(err, "failed to remove marked as deleted paths")
		}
	}

	return nil
}

func mergeSourceIntoTarget(sourceFolder string, targetFolder string, options *mergeoptions.MergeOptions) error {
	osFileSystem := afero.NewOsFs()

	if walkError := afero.Walk(osFileSystem, sourceFolder, func(path string, info fs.FileInfo, _ error) error {
		return mergePathIntoTarget(path, info, sourceFolder, options, targetFolder, osFileSystem)
	}); walkError != nil {
		return errorx.ExternalError.Wrap(walkError, "Failed to walk through source folder")
	}

	if !options.DryRun && sourceFolder != targetFolder && !options.CopyMode {
		if err := pathutils.CleanupPath(sourceFolder); err != nil {
			return err
		}
	}

	return nil
}

func mergePathIntoTarget(
	path string,
	info fs.FileInfo,
	sourceFolder string,
	options *mergeoptions.MergeOptions,
	targetFolder string,
	osFileSystem afero.Fs,
) error {
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
		if err := foldermover.MoveFolder(path, sourceFolder, targetFolder, osFileSystem, options); err != nil {
			return errorx.InternalError.Wrap(err, "Failed to move folder")
		}
	} else {
		if err := filemover.MoveFile(path, sourceFolder, targetFolder, osFileSystem, options); err != nil {
			return errorx.InternalError.Wrap(err, "Failed to move file")
		}
	}

	return nil
}
