package overlay

import (
	"fmt"
	"github.com/containers/storage/pkg/reexec"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
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

	nsContext := NamespaceContext{}
	nsContext.convertFromStringArgs(os.Args[1:])

	var isolator Isolate = newChangeIsolatorOverlayfsMergerfsStrategy(
		newChangeIsolator(nsContext.RootFolder, nsContext.ChangeCaptureFolder, nsContext.OperationDirectory, nsContext.WorkingDirectory),
	)

	if err := isolator.setupFolders(); err != nil {
		log.Fatalf("Error setting up defined folders - %s", err)
	}

	isolator.initialize()

	if err := isolator.prepare(); err != nil {
		log.Fatalf("Error in preparing isolation - %s", err)
	}

	if err := nsRun(nsContext); err != nil {
		log.Fatalf("Error in running isolated process - %s", err)
	}

	if err := isolator.cleanup(); err != nil {
		log.Fatalf("Error in preparing isolation - %s", err)
	}
}

func nsRun(nsContext NamespaceContext) error {
	commandString := strings.Join(nsContext.Command, " ")
	log.Debugf("Running command \"%s\" in isolated environment", commandString)
	cmd := exec.Command("/bin/bash", "-c", commandString)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{"PS1=-[chast-ns-process]- # "}

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "Error running command")
	}
	return nil
}

func launchProcessInNewUserNamespace(nsContext ArgsConverter) error {
	var args = append([]string{"nsExecution"}, nsContext.toStringArgs()...)
	cmd := reexec.Command(args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = buildSysProcAttr(false)

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
	if err := checkIfFolderExists(context.RootFolder); err != nil {
		return err
	}

	return launchProcessInNewUserNamespace(context)
}

func checkIfFolderExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.Wrap(err, fmt.Sprintf("Folder %s does not exist", path))
	}
	return nil
}
