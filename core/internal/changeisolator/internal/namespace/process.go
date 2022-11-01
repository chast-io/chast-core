package namespace

import (
	"chast.io/core/internal/changeisolator/pkg/namespace"
	"encoding/json"
	"github.com/containers/storage/pkg/reexec"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"strings"
)

func init() {
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

	isolator, err := nsContext.BuildIsolationStrategy()
	if err != nil {
		log.Fatalf("Cannot load isolation strategy: %v", err)
	}

	if err := isolator.PrepareInsideNS(); err != nil {
		log.Fatalf("Error in preparing isolation - %s", err)
	}

	if err := nsRun(nsContext); err != nil {
		log.Fatalf("Error in running isolated process - %s", err)
	}

	if err := isolator.CleanupInsideNS(); err != nil {
		log.Fatalf("Error in cleaning up isolation - %s", err)
	}
}

func loadNamespaceContext() namespace.Context {
	nsContext := namespace.Context{}
	pipe := os.NewFile(uintptr(3), "pipe")
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

func nsRun(nsContext namespace.Context) error {
	for _, command := range nsContext.Commands {
		commandString := strings.Join(command, " ")
		log.Debugf("Running command \"%s\" in isolated environment", commandString)
		cmd := exec.Command("/bin/bash", "-c", commandString) // TODO make runner configurable

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.Env = []string{"PS1=-[chast-ns-process]- # "}

		if err := cmd.Run(); err != nil {
			return errors.Wrap(err, "Error running command")
		}
	}
	return nil
}
