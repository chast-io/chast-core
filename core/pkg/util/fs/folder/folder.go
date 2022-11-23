package folder

import (
	chastlog "chast.io/core/internal/logger"
	"github.com/joomcode/errorx"
	"github.com/spf13/afero"
)

func IsFolderEmpty(path string) (bool, error) {
	osFileSystem := afero.NewOsFs()

	chastlog.Log.Trace("Checking if folder is empty: ", path)
	empty, err := afero.IsEmpty(osFileSystem, path)

	if err != nil {
		chastlog.Log.Errorf("Error checking if folder %s is empty: %s", path, err)
		return empty, errorx.ExternalError.Wrap(err, "Failed to check if folder is empty")
	}

	return empty, nil
}

func DoesFolderExist(path string) bool {
	osFileSystem := afero.NewOsFs()
	exists, _ := afero.Exists(osFileSystem, path)

	return exists
}
