package handler

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"os"
)

type ChangeRootHandler struct {
	RootFsPath                 string
	WorkingDirectory           string
	originalRootFileDescriptor *os.File
}

func NewChangeRoot(rootFsPath string, workingDirectory string) *ChangeRootHandler {
	return &ChangeRootHandler{
		RootFsPath:       rootFsPath,
		WorkingDirectory: workingDirectory,
	}
}

func (crh *ChangeRootHandler) Init() error {
	root, err := os.Open("/")
	if err != nil {
		return err
	}
	crh.originalRootFileDescriptor = root

	return nil
}

func (crh *ChangeRootHandler) Open() error {
	if err := unix.Chroot(crh.RootFsPath); err != nil {
		if err := crh.Close(); err != nil {
			return err
		}
		return err
	}
	log.Tracef("Successfully changed root to %s", crh.RootFsPath)
	log.Tracef("Trying to change working directory to %s", crh.WorkingDirectory)
	if err := unix.Chdir(crh.WorkingDirectory); err != nil {
		return err
	}

	return nil
}

func (crh *ChangeRootHandler) Close() error {
	if err := crh.originalRootFileDescriptor.Chdir(); err != nil {
		return err
	}

	if err := unix.Chroot("."); err != nil {
		return err
	}

	if err := crh.originalRootFileDescriptor.Close(); err != nil {
		return err
	}
	return nil
}
