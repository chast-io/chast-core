package fs

import (
	"io"
	"os"
)

func IsFolderEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}

	if err := f.Close(); err != nil {
		return false, err
	}
	return false, err
}

func DoesFolderExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
