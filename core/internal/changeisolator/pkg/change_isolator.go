package pkg

import (
	namespaceInternal "chast.io/core/internal/changeisolator/internal/namespace"
	"chast.io/core/internal/changeisolator/pkg/namespace"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func RunCommandInIsolatedEnvironment(nsContext *namespace.Context) error {
	log.SetLevel(log.TraceLevel)

	userNamespaceRunnerContext := namespaceInternal.New(nsContext)
	if err := userNamespaceRunnerContext.Initialize(); err != nil {
		return errors.Wrap(err, "Error initializing user namespace runner context")
	}

	if err := userNamespaceRunnerContext.Run(); err != nil {
		return errors.Wrap(err, "Failed to run command in isolated environment")
	}

	return nil
}
