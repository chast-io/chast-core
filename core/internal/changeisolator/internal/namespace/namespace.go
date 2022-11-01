package namespace

import (
	"bytes"
	"chast.io/core/internal/changeisolator/internal/strategie"
	"chast.io/core/internal/changeisolator/pkg/namespace"
	"chast.io/core/pkg/util/fs"
	"encoding/json"
	"fmt"
	"github.com/containers/storage/pkg/reexec"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"syscall"
)

type UserNamespaceRunnerContext struct {
	nsContext *namespace.Context
	isolator  strategie.Isolator
}

func New(nsContext *namespace.Context) *UserNamespaceRunnerContext {
	return &UserNamespaceRunnerContext{
		nsContext: nsContext,
	}
}

func (nsrc *UserNamespaceRunnerContext) Initialize() error {
	if !fs.DoesFolderExist(nsrc.nsContext.RootFolder) {
		return errors.Errorf("Root folder %s does not exist", nsrc.nsContext.RootFolder)
	}

	isolator, err := nsrc.nsContext.BuildIsolationStrategy()
	if err != nil {
		return err
	}
	nsrc.isolator = isolator

	return nil
}

func (nsrc *UserNamespaceRunnerContext) Run() error {
	isolator := nsrc.isolator

	if err := isolator.Initialize(); err != nil {
		return err
	}

	if err := isolator.PrepareOutsideNS(); err != nil {
		return err
	}

	if err := nsrc.launchProcess(); err != nil {
		return err
	}

	if err := isolator.CleanupOutsideNS(); err != nil {
		return err
	}

	return nil
}

func (nsrc *UserNamespaceRunnerContext) launchProcess() error {
	cmd, err := nsrc.setupCommand()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return errors.Errorf("Error starting the reexec.Command - %s\n", err)
	}

	if err := cmd.Wait(); err != nil {
		return errors.Errorf("Error waiting for the reexec.Command - %s\n", err)
	}

	return nil
}

func (nsrc *UserNamespaceRunnerContext) setupCommand() (*exec.Cmd, error) {
	nsContext := nsrc.nsContext

	cmd := reexec.Command(processExecutionFunction)

	// https://github.com/containers/buildah/blob/main/run_common.go#L1097
	// setPdeathsig(cmd)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = buildSysProcAttr(false)

	// https://github.com/containers/buildah/blob/main/run_common.go#L1126
	//cmd.Env = util.MergeEnv(os.Environ(), []string{fmt.Sprintf("LOGLEVEL=%d", log.GetLevel())})

	namespaceContextFile, err := buildNamespaceContextFile(nsContext)
	if err != nil {
		return nil, err
	}

	cmd.ExtraFiles = append([]*os.File{namespaceContextFile}, cmd.ExtraFiles...)

	return cmd, nil
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

func buildNamespaceContextFile(nsContext *namespace.Context) (*os.File, error) {
	encodedNsContext, marshalingErr := json.Marshal(nsContext)
	if marshalingErr != nil {
		return nil, fmt.Errorf("encoding configuration for %q: %w", nsContext, marshalingErr)
	}

	pipeReader, pipeWriter, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("creating configuration pipe: %w", err)
	}
	_, nsContextCopyError := io.Copy(pipeWriter, bytes.NewReader(encodedNsContext))
	if nsContextCopyError != nil {
		nsContextCopyError = fmt.Errorf("while copying configuration down pipe to child process: %w", nsContextCopyError)
	}

	if err := pipeWriter.Close(); err != nil {
		return nil, errors.Wrap(err, "Error closing config pipe writer")
	}

	return pipeReader, nil
}
