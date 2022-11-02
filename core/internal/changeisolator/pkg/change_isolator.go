package pkg

import (
	"chast.io/core/internal/changeisolator/internal/namespace"
	namespace2 "chast.io/core/internal/changeisolator/pkg/namespace"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func RunCommandInIsolatedEnvironment(nsContext *namespace2.Context) error {
	log.SetLevel(log.TraceLevel)

	userNamespaceRunnerContext := namespace.New(nsContext)
	if err := userNamespaceRunnerContext.Initialize(); err != nil {
		return errors.Wrap(err, "Error initializing user namespace runner context")
	}

	if err := userNamespaceRunnerContext.Run(); err != nil {
		return errors.Wrap(err, "Failed to run command in isolated environment")
	}

	return nil
}
