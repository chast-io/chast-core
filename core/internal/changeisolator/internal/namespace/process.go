package namespace

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strings"

	"chast.io/core/internal/changeisolator/pkg/namespace"
	"github.com/containers/storage/pkg/reexec"
	log "github.com/sirupsen/logrus"
)

func init() { //nolint:gochecknoinits // This function needs to register the function used for reexec
	reexec.Register(processExecutionFunction, nsExecution)

	if reexec.Init() {
		os.Exit(0)
	}
}

const processExecutionFunction = "nsExecution"

func nsExecution() {
	log.SetLevel(log.TraceLevel)

	log.Printf("Running in isolated environment")

	nsContext := loadNamespaceContext()

	isolator, isolationStrategyBuildError := nsContext.BuildIsolationStrategy()
	if isolationStrategyBuildError != nil {
		log.Fatalf("Cannot load isolation strategy: %v", isolationStrategyBuildError)
	}

	if err := isolator.PrepareInsideNS(); err != nil {
		log.Fatalf("Error in preparing isolation - %s", err)
	}

	nsRun(nsContext)

	if err := isolator.CleanupInsideNS(); err != nil {
		log.Fatalf("Error in cleaning up isolation - %s", err)
	}
}

const firstExtraFileFileDescriptorNumber = 3

func loadNamespaceContext() namespace.Context {
	nsContext := *namespace.NewEmptyContext()
	pipe := os.NewFile(uintptr(firstExtraFileFileDescriptorNumber), "pipe")

	data, err := io.ReadAll(pipe)
	if err != nil {
		log.Fatalf("Error while reading namespace context from pipe: %v", err)
	}

	err = json.Unmarshal(data, &nsContext)
	if err != nil {
		log.Fatalf("Error while decoding namespace context: %v", err)
	}

	return nsContext
}

func nsRun(nsContext namespace.Context) {
	for _, command := range nsContext.Commands {
		commandString := strings.Join(command, " ")
		log.Debugf("Running command \"%s\" in isolated environment", commandString)

		cmd := exec.Command("/bin/bash", "-c", commandString) // TODO make runner configurable

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.Env = []string{"PS1=-[chast-ns-process]- # "}

		if err := cmd.Run(); err != nil {
			log.Warnf("Error running command: %v", err)
		}
	}
}