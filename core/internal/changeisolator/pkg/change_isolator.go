package changeisolator

import (
	namespaceInternal "chast.io/core/internal/changeisolator/internal/namespace"
	"chast.io/core/internal/changeisolator/pkg/namespace"
	"github.com/joomcode/errorx"
)

func RunCommandInIsolatedEnvironment(nsContext *namespace.Context) error {
	userNamespaceRunnerContext := namespaceInternal.New(nsContext)
	if err := userNamespaceRunnerContext.Initialize(); err != nil {
		return errorx.InternalError.Wrap(err, "Error initializing user namespace runner context")
	}

	if err := userNamespaceRunnerContext.Run(); err != nil {
		return errorx.InternalError.Wrap(err, "Failed to run command in isolated environment")
	}

	return nil
}
