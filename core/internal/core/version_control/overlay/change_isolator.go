package overlay

import (
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

type changeIsolator struct {
	Isolate

	RootFolder          string
	ChangeCaptureFolder string
	OperationDirectory  string
	WorkingDirectory    string
}

type Isolate interface {
	initialize() error
	prepareOutsideNS() error
	prepareInsideNS() error
	cleanupOutsideNS() error
	cleanupInsideNS() error
	getIsolationStrategy() IsolationStrategy
}

func newChangeIsolator(
	rootFolder string,
	changeCaptureFolder string,
	operationDirectory string,
	currentWorkingDirectory string) *changeIsolator {
	return &changeIsolator{
		RootFolder:          rootFolder,
		ChangeCaptureFolder: changeCaptureFolder,
		OperationDirectory:  operationDirectory,
		WorkingDirectory:    currentWorkingDirectory,
	}
}

func (changeIsolator *changeIsolator) initialize() error {
	return changeIsolator.setupFolders()
}

func (changeIsolator *changeIsolator) setupFolders() error {
	log.Printf("Setting up folders: %s, %s, \n", changeIsolator.ChangeCaptureFolder, changeIsolator.OperationDirectory)
	if err := os.MkdirAll(changeIsolator.ChangeCaptureFolder, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(changeIsolator.OperationDirectory, 0755); err != nil {
		return err
	}
	return nil
}

func (changeIsolator *changeIsolator) cleanupInsideNS() error {
	log.Printf("[Inside NS ] Cleaning up change isolator")
	return nil
}

func (changeIsolator *changeIsolator) cleanupOutsideNS() error {
	log.Printf("[Outside NS ] Cleaning up change isolator")

	isEmpty, err := isFolderEmpty(changeIsolator.OperationDirectory)
	if err != nil {
		return err
	}
	if isEmpty {
		err := os.RemoveAll(changeIsolator.OperationDirectory)
		if err != nil {
			return err
		}
	}
	return nil
}

func isFolderEmpty(name string) (bool, error) {
	f, err := os.Open(name)
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
