package folder

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

func IsFolderEmpty(path string) (bool, error) {
	osFileSystem := afero.NewOsFs()
	empty, err := afero.IsEmpty(osFileSystem, path)

	return empty, errors.Wrap(err, "Failed to check if folder is empty")
}

func DoesFolderExist(path string) bool {
	osFileSystem := afero.NewOsFs()
	exists, _ := afero.Exists(osFileSystem, path)

	return exists
}
