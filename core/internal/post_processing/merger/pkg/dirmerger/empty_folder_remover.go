package dirmerger

import (
	"io/fs"
	"os"
	"path/filepath"

	chastlog "chast.io/core/internal/logger"
	"github.com/joomcode/errorx"
	"github.com/spf13/afero"
)

func removeEmptyFolders(
	folder string,
	options *MergeOptions,
) error {
	if options.DryRun {
		return nil
	}

	osFileSystem := afero.NewOsFs()

	if walkError := afero.Walk(osFileSystem, folder, func(path string, info fs.FileInfo, _ error) error {
		if info == nil {
			return nil
		}

		if info.IsDir() {
			if err := removeFolderAndParentsIfEmpty(path); err != nil {
				return err
			}
		}

		return nil
	}); walkError != nil {
		return errorx.ExternalError.Wrap(walkError, "Failed to walk through target folder")
	}

	return nil
}

func removeFolderAndParentsIfEmpty(path string) error {
	osFileSystem := afero.NewOsFs()

	exists, existsCheckErr := afero.Exists(osFileSystem, path)
	if existsCheckErr != nil {
		return errorx.ExternalError.Wrap(existsCheckErr, "Failed to check if path exists")
	}

	if !exists {
		return nil
	}

	isDir, isDirCheckError := afero.IsDir(osFileSystem, path)
	if isDirCheckError != nil {
		return errorx.ExternalError.Wrap(isDirCheckError, "Failed to check if path is a folder")
	}

	if isDir {
		isEmpty, isEmptyError := afero.IsEmpty(osFileSystem, path)
		if isEmptyError != nil {
			return errorx.InternalError.Wrap(isEmptyError, "Failed to check if folder is empty")
		}

		if !isEmpty {
			return nil
		}

		chastlog.Log.Debugf("Empty folder found during merge of folders. Removing: %s", path)

		if removeError := osFileSystem.Remove(path); removeError != nil {
			return errorx.ExternalError.Wrap(removeError, "Failed to remove folder")
		}
	}

	if path == string(os.PathSeparator) {
		return nil
	}

	return removeFolderAndParentsIfEmpty(filepath.Dir(path))
}
