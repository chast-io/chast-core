package overlay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/containers/storage/pkg/reexec"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func init() {
	reexec.Register("nsExecution", nsExecution)
	if reexec.Init() {
		os.Exit(0)
	}
}

func nsExecution() {
	log.SetLevel(log.TraceLevel)

	log.Printf("Running in isolated environment")

	nsContext := NamespaceContext{}

	if len(os.Args) < 1 {
		log.Fatal("Missing arguments")
	}

	isolationStrategyStr, err := strconv.ParseUint(os.Args[1], 10, 8)
	if err != nil {
		log.Fatalf("Error parsing isolation strategy of value \"%s\" - %s", os.Args[1], err)
	}

	pipe := os.NewFile(uintptr(3), "pipe")
	data, err := io.ReadAll(pipe)
	if err != nil {
		log.Fatalf("Error while reading namespace context from pipe: %v", err)
	}
	err = json.Unmarshal(data, &nsContext)
	if err != nil {
		log.Fatalf("Error while decoding namespace context: %v", err)
	}

	var isolationStrategy = uint8(isolationStrategyStr)
	isolator, err := nsContext.GetIsolationStrategy(
		isolationStrategy,
		*newChangeIsolator(nsContext.RootFolder, nsContext.ChangeCaptureFolder, nsContext.OperationDirectory, nsContext.WorkingDirectory),
	)
	if err != nil {
		log.Fatalf("Cannot load isolation strategy: %v", err)
	}

	if err := isolator.prepareInsideNS(); err != nil {
		log.Fatalf("Error in preparing isolation - %s", err)
	}

	if err := nsRun(nsContext); err != nil {
		log.Fatalf("Error in running isolated process - %s", err)
	}

	if err := isolator.cleanupInsideNS(); err != nil {
		log.Fatalf("Error in cleaning up isolation - %s", err)
	}
}

func nsRun(nsContext NamespaceContext) error {
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

func launchProcessInNewUserNamespace(nsContext *NamespaceContext, isolationStrategy IsolationStrategy) error {
	var args = append([]string{"nsExecution", strconv.Itoa(int(isolationStrategy))})
	cmd := reexec.Command(args...)

	encodedNsContext, marshalingErr := json.Marshal(nsContext)
	if marshalingErr != nil {
		return fmt.Errorf("encoding configuration for %q: %w", nsContext, marshalingErr)
	}

	// https://github.com/containers/buildah/blob/main/run_common.go#L1097
	// setPdeathsig(cmd)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = buildSysProcAttr(false)

	// https://github.com/containers/buildah/blob/main/run_common.go#L1126
	//cmd.Env = util.MergeEnv(os.Environ(), []string{fmt.Sprintf("LOGLEVEL=%d", log.GetLevel())})

	preader, pwriter, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("creating configuration pipe: %w", err)
	}
	_, nsContextCopyError := io.Copy(pwriter, bytes.NewReader(encodedNsContext))
	if nsContextCopyError != nil {
		nsContextCopyError = fmt.Errorf("while copying configuration down pipe to child process: %w", nsContextCopyError)
	}

	if err := pwriter.Close(); err != nil {
		return errors.Wrap(err, "Error closing config pipe writer")
	}

	cmd.ExtraFiles = append([]*os.File{preader}, cmd.ExtraFiles...)

	if err := cmd.Start(); err != nil {
		return errors.Errorf("Error starting the reexec.Command - %s\n", err)
	}

	if err := cmd.Wait(); err != nil {
		return errors.Errorf("Error waiting for the reexec.Command - %s\n", err)
	}

	return nil
}

func buildSysProcAttr(networkCapabilitiesRequired bool) *syscall.SysProcAttr {
	userToRootUidMappings := []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID:      os.Getuid(),
			Size:        1,
		},
	}

	userGroupToRootGroupMappings := []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID:      os.Getgid(),
			Size:        1,
		},
	}

	var cloneFlags uintptr = syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER

	if !networkCapabilitiesRequired {
		cloneFlags |= syscall.CLONE_NEWNET
	}

	return &syscall.SysProcAttr{
		UidMappings: userToRootUidMappings,
		GidMappings: userGroupToRootGroupMappings,
		Cloneflags:  cloneFlags,
	}
}

func RunCommandInIsolatedEnvironment(context *NamespaceContext) error {
	log.SetLevel(log.TraceLevel)

	if err := checkIfFolderExists(context.RootFolder); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Root folder %s does not exist", context.RootFolder))
	}

	//var isolator Isolate = newChangeIsolatorOverlayfsMergerfsStrategy(
	//	*newChangeIsolator(context.RootFolder, context.ChangeCaptureFolder, context.OperationDirectory, context.WorkingDirectory),
	//)

	var isolator Isolate = newChangeIsolatorUnionFsStrategy(
		*newChangeIsolator(context.RootFolder, context.ChangeCaptureFolder, context.OperationDirectory, context.WorkingDirectory),
	)

	if err := isolator.initialize(); err != nil {
		return err
	}

	if err := isolator.prepareOutsideNS(); err != nil {
		return err
	}

	if err := launchProcessInNewUserNamespace(context, isolator.getIsolationStrategy()); err != nil {
		return err
	}

	if err := isolator.cleanupOutsideNS(); err != nil {
		return err
	}

	return nil
}

func checkIfFolderExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.Wrap(err, fmt.Sprintf("Folder %s does not exist", path))
	}
	return nil
}
