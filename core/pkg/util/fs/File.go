package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

type File struct {
	FileReader

	FileName        string
	ParentDirectory string
	AbsolutePath    string
	path            string
	data            *[]byte
}

type FileReader interface {
	Read() *[]byte
}

func NewFile(path string) *File {
	file := File{}
	file.path = path

	dirName, fileName := filepath.Split(path)
	file.FileName = fileName
	file.ParentDirectory = dirName

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		panic(fmt.Sprintf("Could not load current directory"))
	}
	file.AbsolutePath = absolutePath

	return &file
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
