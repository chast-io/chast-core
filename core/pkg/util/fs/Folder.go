package fs

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

func IsFolderEmpty(path string) (bool, error) {
	folder, openDirError := os.Open(path)
	if openDirError != nil {
		return false, errors.Wrap(openDirError, "Could not open folder")
	}

	_, readDirErr := folder.Readdirnames(1)
	if errors.Is(readDirErr, io.EOF) {
		return true, nil
	}

	if closeErr := folder.Close(); closeErr != nil {
		return false, errors.Wrap(closeErr, "Failed to close folder")
	}

	return false, errors.Wrap(readDirErr, "Error while reading directory")
}

func DoesFolderExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
