package pkg

import (
	"chast.io/core/internal/changeisolator/internal/namespace"
	namespace2 "chast.io/core/internal/changeisolator/pkg/namespace"
	log "github.com/sirupsen/logrus"
)

func RunCommandInIsolatedEnvironment(nsContext *namespace2.Context) error {
	log.SetLevel(log.TraceLevel)

	userNamespaceRunnerContext := namespace.New(nsContext)
	if err := userNamespaceRunnerContext.Initialize(); err != nil {
		return err
	}

	return userNamespaceRunnerContext.Run()
}
