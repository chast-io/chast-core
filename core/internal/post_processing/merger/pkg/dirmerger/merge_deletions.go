package dirmerger

import (
	"strings"

	"chast.io/core/internal/post_processing/merger/pkg/mergeoptions"
	"chast.io/core/internal/post_processing/merger/pkg/mergererrors"
	"github.com/joomcode/errorx"
	"github.com/spf13/afero"
)

func RemoveMarkedAsDeletedPaths(locations []string, options *mergeoptions.MergeOptions) error {
	osFileSystem := afero.NewOsFs()

	for _, location := range locations {
		targetExists, targetExistsError := afero.Exists(osFileSystem, location)
		if targetExistsError != nil {
			return errorx.ExternalError.Wrap(targetExistsError, "Failed to check if target folder exists")
		}

		if !targetExists {
			return nil
		}

		if strings.HasSuffix(location, options.MetaFilesDeletedExtension) {
			mergeDeletedPathError := removeMarkedAsDeletedPath(location, osFileSystem, options)
			if mergeDeletedPathError != nil {
				return errorx.InternalError.Wrap(mergeDeletedPathError, "Failed to merge deleted path")
			}
		}
	}

	return nil
}

func removeMarkedAsDeletedPath(
	path string,
	osFileSystem afero.Fs,
	options *mergeoptions.MergeOptions,
) error {
	exists, existsError := afero.Exists(osFileSystem, path)
	if existsError != nil {
		return errorx.ExternalError.Wrap(existsError, "Failed to check if path exists")
	}

	if exists {
		if options.BlockOverwrite {
			return errorx.InternalError.Wrap(mergererrors.ErrMergeOverwriteBlock,
				"Failed to remove marked as deleted path: %s", path)
		}

		if !options.DryRun {
			if err := osFileSystem.RemoveAll(path); err != nil {
				return errorx.ExternalError.Wrap(err, "Failed to remove file")
			}
		}
	}

	return nil
}
