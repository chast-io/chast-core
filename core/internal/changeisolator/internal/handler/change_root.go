package handler

import (
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

type ChangeRootHandler struct {
	RootFsPath                 string
	WorkingDirectory           string
	originalRootFileDescriptor *os.File
}

func NewChangeRoot(rootFsPath string, workingDirectory string) *ChangeRootHandler {
	return &ChangeRootHandler{
		RootFsPath:                 rootFsPath,
		WorkingDirectory:           workingDirectory,
		originalRootFileDescriptor: nil,
	}
}

func (crh *ChangeRootHandler) Init() error {
	root, pathOpenError := os.Open("/")
	if pathOpenError != nil {
		return errors.Wrap(pathOpenError, "Failed to open root path")
	}

	crh.originalRootFileDescriptor = root

	return nil
}

func (crh *ChangeRootHandler) Open() error {
	if chrootErr := unix.Chroot(crh.RootFsPath); chrootErr != nil {
		if closeErr := crh.Close(); closeErr != nil {
			return errors.Wrap(closeErr, "Failed to close change root handler")
		}

		return errors.Wrap(chrootErr, "Failed to change root")
	}

	log.Tracef("Successfully changed root to %s", crh.RootFsPath)
	log.Tracef("Trying to change working directory to %s", crh.WorkingDirectory)

	if err := unix.Chdir(crh.WorkingDirectory); err != nil {
		return errors.Wrap(err, "Failed to change working directory")
	}

	return nil
}

func (crh *ChangeRootHandler) Close() error {
	if err := crh.originalRootFileDescriptor.Chdir(); err != nil {
		return errors.Wrap(err, "Failed to change directory to original root")
	}

	if err := unix.Chroot("."); err != nil {
		return errors.Wrap(err, "Failed to change root to original root")
	}

	if err := crh.originalRootFileDescriptor.Close(); err != nil {
		return errors.Wrap(err, "Failed to close original root file descriptor")
	}

	return nil
}
