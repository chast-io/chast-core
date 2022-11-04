package strategie

import (
	"os"

	"chast.io/core/internal/changeisolator/pkg/strategy"
	"chast.io/core/pkg/util/fs/folder"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
	currentWorkingDirectory string,
) *IsolatorContext {
	return &IsolatorContext{
		RootFolder:          rootFolder,
		ChangeCaptureFolder: changeCaptureFolder,
		OperationDirectory:  operationDirectory,
		WorkingDirectory:    currentWorkingDirectory,

		Isolator: nil,
	}
}

func (changeIsolator *IsolatorContext) Initialize() error {
	return changeIsolator.setupFolders()
}

func (changeIsolator *IsolatorContext) setupFolders() error {
	log.Printf("Setting up folders: %s, %s, \n", changeIsolator.ChangeCaptureFolder, changeIsolator.OperationDirectory)

	if err := os.MkdirAll(changeIsolator.ChangeCaptureFolder, 0o755); err != nil {
		return errors.Wrap(err, "Error creating change capture folder")
	}

	if err := os.MkdirAll(changeIsolator.OperationDirectory, 0o755); err != nil {
		return errors.Wrap(err, "Error creating operation folder")
	}

	return nil
}

func (changeIsolator *IsolatorContext) CleanupInsideNS() error {
	log.Tracef("[Inside NS] Cleaning up change isolator")

	return nil
}

func (changeIsolator *IsolatorContext) CleanupOutsideNS() error {
	log.Tracef("[Outside NS] Cleaning up change isolator")

	isEmpty, isFolderEmptyError := folder.IsFolderEmpty(changeIsolator.OperationDirectory)
	if isFolderEmptyError != nil {
		return errors.Wrap(isFolderEmptyError, "Error checking if operation directory is empty")
	}

	if isEmpty {
		if err := os.RemoveAll(changeIsolator.OperationDirectory); err != nil {
			return errors.Wrap(err, "Error removing operation directory")
		}
	}

	return nil
}
