package change_isolator

import (
	"github.com/pkg/errors"
)

type NamespaceContext struct {
	RootFolder          string
	MergeFolders        []string
	ChangeCaptureFolder string
	OperationDirectory  string
	WorkingDirectory    string
	Commands            [][]string
}

type IsolationStrategy = uint8

const (
	UnknownIsolation IsolationStrategy = iota
	OverlayFS        IsolationStrategy = iota
	UnionFS          IsolationStrategy = iota
)

func NewNamespaceContext(
	rootFolder string,
	mergeFolders []string,
	changeCaptureFolder string,
	operationDirectory string,
	workingDirectory string,
	command [][]string) *NamespaceContext {
	return &NamespaceContext{
		RootFolder:          rootFolder,
		MergeFolders:        mergeFolders,
		ChangeCaptureFolder: changeCaptureFolder,
		OperationDirectory:  operationDirectory,
		WorkingDirectory:    workingDirectory,
		Commands:            command,
	}
}

func (nsc *NamespaceContext) GetIsolationStrategy(strategy IsolationStrategy, changeIsolator changeIsolator) (Isolate, error) {
	var isolator Isolate
	switch strategy {
	case OverlayFS:
		isolator = newChangeIsolatorOverlayfsMergerfsStrategy(changeIsolator)
	case UnionFS:
		isolator = newChangeIsolatorUnionFsStrategy(changeIsolator)
	default:
		return nil, errors.Errorf("unknown isolation strategy")
	}

	if err := isolator.initialize(); err != nil {
		return nil, err
	}
	return isolator, nil
}
