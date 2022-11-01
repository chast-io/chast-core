package strategie

import (
	"chast.io/core/internal/changeisolator/pkg/strategy"
	"chast.io/core/pkg/util/fs"
	log "github.com/sirupsen/logrus"
	"os"
)

type Isolator interface {
	Initialize() error
	PrepareOutsideNS() error
	PrepareInsideNS() error
	CleanupOutsideNS() error
	CleanupInsideNS() error
	GetIsolationStrategy() strategy.IsolationStrategy
}

type IsolatorContext struct {
	Isolator

	RootFolder          string
	ChangeCaptureFolder string
	OperationDirectory  string
	WorkingDirectory    string
}

func NewChangeIsolator(
	rootFolder string,
	changeCaptureFolder string,
	operationDirectory string,
	currentWorkingDirectory string) *IsolatorContext {
	return &IsolatorContext{
		RootFolder:          rootFolder,
		ChangeCaptureFolder: changeCaptureFolder,
		OperationDirectory:  operationDirectory,
		WorkingDirectory:    currentWorkingDirectory,
	}
}

func (changeIsolator *IsolatorContext) Initialize() error {
	return changeIsolator.setupFolders()
}

func (changeIsolator *IsolatorContext) setupFolders() error {
	log.Printf("Setting up folders: %s, %s, \n", changeIsolator.ChangeCaptureFolder, changeIsolator.OperationDirectory)
	if err := os.MkdirAll(changeIsolator.ChangeCaptureFolder, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(changeIsolator.OperationDirectory, 0755); err != nil {
		return err
	}
	return nil
}

func (changeIsolator *IsolatorContext) CleanupInsideNS() error {
	log.Tracef("[Inside NS] Cleaning up change isolator")
	return nil
}

func (changeIsolator *IsolatorContext) CleanupOutsideNS() error {
	log.Tracef("[Outside NS] Cleaning up change isolator")

	isEmpty, err := fs.IsFolderEmpty(changeIsolator.OperationDirectory)
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
