package dirmerger

import (
	"io/fs"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

var errLocationDoesNotExist = errors.New("Location does not exist")

func MergeDeletions(targetFolder string) error {
	osFileSystem := afero.NewOsFs()

	targetExists, targetExistsError := afero.Exists(osFileSystem, targetFolder)
	if targetExistsError != nil {
		return errors.Wrap(targetExistsError, "Failed to check if target folder exists")
	}

	if !targetExists {
		return errLocationDoesNotExist
	}

	if walkError := afero.Walk(osFileSystem, targetFolder, func(path string, info fs.FileInfo, _ error) error {
		if info == nil {
			return nil
		}

		mergeDeletedPathError := mergeDeletedPath(path, osFileSystem)
		if mergeDeletedPathError != nil {
			return errors.Wrap(mergeDeletedPathError, "Failed to merge deleted path")
		}

		return nil
	}); walkError != nil {
		return errors.Wrap(walkError, "Failed to walk through target folder")
	}

	return nil
}

func mergeDeletedPath(
	path string,
	osFileSystem afero.Fs,
) error {
	if !strings.HasSuffix(path, unionFsHiddenPathSuffix) {
		return nil
	}

	undeletedTargetPath := strings.TrimSuffix(path, unionFsHiddenPathSuffix)

	exists, existsError := afero.Exists(osFileSystem, undeletedTargetPath)
	if existsError != nil {
		return errors.Wrap(existsError, "Failed to check if path exists")
	}

	if exists {
		if err := osFileSystem.RemoveAll(undeletedTargetPath); err != nil {
			return errors.Wrap(err, "Failed to remove file")
		}
	}

	return nil
}
