package folder

import (
	"github.com/joomcode/errorx"
	"github.com/spf13/afero"
)

func IsFolderEmpty(path string) (bool, error) {
	osFileSystem := afero.NewOsFs()
	empty, err := afero.IsEmpty(osFileSystem, path)

	return empty, errorx.ExternalError.Wrap(err, "Failed to check if folder is empty")
}

func DoesFolderExist(path string) bool {
	osFileSystem := afero.NewOsFs()
	exists, _ := afero.Exists(osFileSystem, path)

	return exists
}
