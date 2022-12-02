package namespace

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strings"

	"chast.io/core/internal/changeisolator/pkg/namespace"
	chastlog "chast.io/core/internal/logger"
	"github.com/containers/storage/pkg/reexec"
	"github.com/ttacon/chalk"
)

func init() { //nolint:gochecknoinits // This function needs to register the function used for reexec
	reexec.Register(processExecutionFunction, nsExecution)

	if reexec.Init() {
		os.Exit(0)
	}
}

const processExecutionFunction = "nsExecution"

func nsExecution() {
	chastlog.Log.Printf("Running in isolated environment")

	nsContext := loadNamespaceContext()

	isolator, isolationStrategyBuildError := nsContext.BuildIsolationStrategy()
	if isolationStrategyBuildError != nil {
		chastlog.Log.Fatalf("Cannot load isolation strategy: %v", isolationStrategyBuildError)
	}

	if err := isolator.PrepareInsideNS(); err != nil {
		chastlog.Log.Fatalf("Error in preparing isolation - %s", err)
	}

	nsRun(nsContext)

	if err := isolator.CleanupInsideNS(); err != nil {
		chastlog.Log.Fatalf("Error in cleaning up isolation - %s", err)
	}
}

const firstExtraFileFileDescriptorNumber = 3

func loadNamespaceContext() namespace.Context {
	nsContext := *namespace.NewEmptyContext()
	pipe := os.NewFile(uintptr(firstExtraFileFileDescriptorNumber), "pipe")

	data, err := io.ReadAll(pipe)
	if err != nil {
		chastlog.Log.Fatalf("Error while reading namespace context from pipe: %v", err)
	}

	err = json.Unmarshal(data, &nsContext)
	if err != nil {
		chastlog.Log.Fatalf("Error while decoding namespace context: %v", err)
	}

	return nsContext
}

func nsRun(nsContext namespace.Context) {
	for _, command := range nsContext.Commands {
		commandString := strings.Join(command, " ")
		chastlog.Log.Debugf("Running command \"%s\" in isolated environment", chalk.Blue.Color(commandString))

		cmd := exec.Command("/bin/bash", "-c", commandString) // TODO make runner configurable

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.Env = []string{"PS1=-[chast-ns-process]- # "}

		if err := cmd.Run(); err != nil {
			chastlog.Log.Warnf("Error running command: %v", err)
		}

		chastlog.Log.Debugf("Running command done!")
	}
}
