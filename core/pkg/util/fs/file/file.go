package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joomcode/errorx"
)

type File struct {
	FileName        string
	ParentDirectory string
	AbsolutePath    string
	path            string
	data            *[]byte
}

func NewFile(path string) (*File, error) {
	dirName, fileName := filepath.Split(path)

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, errorx.ExternalError.Wrap(err, "Could not load current directory")
	}

	return &File{ //nolint:exhaustruct // data is used as cache
		path:            path,
		FileName:        fileName,
		ParentDirectory: dirName,
		AbsolutePath:    absolutePath,
	}, nil
}

func (file *File) Exists() bool {
	_, err := os.Stat(file.path)

	return !os.IsNotExist(err)
}

func (file *File) Read() *[]byte {
	if file.data == nil {
		fileContent, err := os.ReadFile(file.path)
		if err != nil {
			panic(fmt.Sprintf("File '%s' does not exist.", file.path))
		}

		file.data = &fileContent
	}

	return file.data
}
