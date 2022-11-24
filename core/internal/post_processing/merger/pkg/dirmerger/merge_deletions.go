package dirmerger

import (
	"io/fs"
	"strings"

	"github.com/joomcode/errorx"
	"github.com/spf13/afero"
)

func removeMarkedAsDeletedPaths(targetFolder string, options *MergeOptions) error {
	osFileSystem := afero.NewOsFs()

	targetExists, targetExistsError := afero.Exists(osFileSystem, targetFolder)
	if targetExistsError != nil {
		return errorx.ExternalError.Wrap(targetExistsError, "Failed to check if target folder exists")
	}

	if !targetExists {
		return errorx.ExternalError.New("Location does not exist")
	}

	if walkError := afero.Walk(osFileSystem, targetFolder, func(path string, info fs.FileInfo, _ error) error {
		if info == nil {
			return nil
		}

		if strings.HasSuffix(path, options.MetaFilesDeletedExtension) {
			mergeDeletedPathError := removeMarkedAsDeletedPath(path, osFileSystem, options)
			if mergeDeletedPathError != nil {
				return errorx.InternalError.Wrap(mergeDeletedPathError, "Failed to merge deleted path")
			}
		}

		return nil
	}); walkError != nil {
		return errorx.ExternalError.Wrap(walkError, "Failed to walk through target folder")
	}

	return nil
}

func removeMarkedAsDeletedPath(
	path string,
	osFileSystem afero.Fs,
	options *MergeOptions,
) error {
	exists, existsError := afero.Exists(osFileSystem, path)
	if existsError != nil {
		return errorx.ExternalError.Wrap(existsError, "Failed to check if path exists")
	}

	if exists {
		if options.BlockOverwrite {
			return errMergeOverwriteBlock
		}

		if !options.DryRun {
			if err := osFileSystem.RemoveAll(path); err != nil {
				return errorx.ExternalError.Wrap(err, "Failed to remove file")
			}
		}
	}

	return nil
}
