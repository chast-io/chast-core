package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

type File struct {
	FileName        string
	ParentDirectory string
	path            string
	absolutePath    string
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
	file.absolutePath = absolutePath

	return &file
}

func (file File) Read() *[]byte {
	if file.data == nil {
		fileContent, err := os.ReadFile(file.path)
		if err != nil {
			panic(fmt.Sprintf("File '%s' does not exist.", file.path))
		}
		file.data = &fileContent
	}
	return file.data
}
