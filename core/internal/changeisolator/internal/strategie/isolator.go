package strategie

import (
	"os"

	"chast.io/core/internal/changeisolator/pkg/strategy"
	chastlog "chast.io/core/internal/logger"
	"chast.io/core/pkg/util/fs/folder"
	"github.com/joomcode/errorx"
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
	RootJoinFolders     []string
	ChangeCaptureFolder string
	OperationDirectory  string
	WorkingDirectory    string
}

func NewChangeIsolator(
	rootFolder string,
	rootJoinFolders []string,
	changeCaptureFolder string,
	operationDirectory string,
	currentWorkingDirectory string,
) *IsolatorContext {
	return &IsolatorContext{
		RootFolder:          rootFolder,
		RootJoinFolders:     rootJoinFolders,
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
	chastlog.Log.Printf(
		"Setting up folders: %s, %s",
		changeIsolator.ChangeCaptureFolder,
		changeIsolator.OperationDirectory,
	)

	if err := os.MkdirAll(changeIsolator.ChangeCaptureFolder, 0o755); err != nil {
		return errorx.ExternalError.Wrap(err, "Error creating change capture folder")
	}

	if err := os.MkdirAll(changeIsolator.OperationDirectory, 0o755); err != nil {
		return errorx.ExternalError.Wrap(err, "Error creating operation folder")
	}

	return nil
}

func (changeIsolator *IsolatorContext) CleanupInsideNS() error {
	chastlog.Log.Tracef("[Inside NS] Cleaning up change isolator")

	return nil
}

func (changeIsolator *IsolatorContext) CleanupOutsideNS() error {
	chastlog.Log.Tracef("[Outside NS] Cleaning up change isolator")

	isEmpty, isFolderEmptyError := folder.IsFolderEmpty(changeIsolator.OperationDirectory)
	if isFolderEmptyError != nil {
		return errorx.InternalError.Wrap(isFolderEmptyError, "Error checking if operation directory is empty")
	}

	if isEmpty {
		if err := os.RemoveAll(changeIsolator.OperationDirectory); err != nil {
			return errorx.ExternalError.Wrap(err, "Error removing operation directory")
		}
	}

	return nil
}
