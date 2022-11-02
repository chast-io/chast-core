package namespace

import (
	"chast.io/core/internal/changeisolator/internal/strategie"
	"chast.io/core/internal/changeisolator/pkg/strategy"
	"github.com/pkg/errors"
)

type Context struct {
	RootFolder          string
	MergeFolders        []string
	ChangeCaptureFolder string
	OperationDirectory  string
	WorkingDirectory    string
	Commands            [][]string

	IsolationStrategy strategy.IsolationStrategy
}

func NewContext(
	rootFolder string,
	mergeFolders []string,
	changeCaptureFolder string,
	operationDirectory string,
	workingDirectory string,
	command [][]string,
	isolationStrategy strategy.IsolationStrategy,
) *Context {
	return &Context{
		RootFolder:          rootFolder,
		MergeFolders:        mergeFolders,
		ChangeCaptureFolder: changeCaptureFolder,
		OperationDirectory:  operationDirectory,
		WorkingDirectory:    workingDirectory,
		Commands:            command,

		IsolationStrategy: isolationStrategy,
	}
}

func NewEmptyContext() *Context {
	return &Context{} //nolint:exhaustruct // initialized empty here for later full initialization
}

func (nsc *Context) BuildIsolationStrategy() (strategie.Isolator, error) { //nolint:ireturn // Factory method
	var isolator strategie.Isolator

	isolatorContext := nsc.newContextFromNamespaceContext()

	switch nsc.IsolationStrategy {
	case strategy.OverlayFS:
		isolator = strategie.NewOverlayFsMergerFsStrategy(isolatorContext)
	case strategy.UnionFS:
		isolator = strategie.NewUnionFsStrategy(isolatorContext)
	default:
		return nil, errors.Errorf("unknown isolation strategy")
	}

	if err := isolator.Initialize(); err != nil {
		return nil, errors.Wrap(err, "Error initializing isolation strategy")
	}

	return isolator, nil
}

func (nsc *Context) newContextFromNamespaceContext() strategie.IsolatorContext {
	return strategie.IsolatorContext{
		RootFolder:          nsc.RootFolder,
		ChangeCaptureFolder: nsc.ChangeCaptureFolder,
		OperationDirectory:  nsc.OperationDirectory,
		WorkingDirectory:    nsc.WorkingDirectory,

		Isolator: nil,
	}
}
