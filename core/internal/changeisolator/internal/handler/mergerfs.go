package handler

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

type MergerFsHandler struct {
	Source string
	Target string
}

func NewMergerFs(source string, target string) *MergerFsHandler {
	return &MergerFsHandler{
		Source: source,
		Target: target,
	}
}

func (mergerFs *MergerFsHandler) Mount() error {
	log.Tracef("Trying to merge %s into %s", mergerFs.Source, mergerFs.Target)

	if err := os.MkdirAll(mergerFs.Target, 0o755); err != nil {
		return errors.Wrap(err, "Failed to create mergerFs target dir")
	}

	command := "/usr/bin/mergerfs"
	args := []string{
		mergerFs.Source,
		mergerFs.Target,
	}

	if _, err := exec.Command(command, args...).CombinedOutput(); err != nil {
		return errors.Wrap(err, "Failed to mount mergerfs")
	}

	log.Debugf("Mergerfs was successfully mounted -  %s into %s", mergerFs.Source, mergerFs.Target)

	return nil
}

func (mergerFs *MergerFsHandler) Unmount() error {
	log.Tracef("Trying to unmerge mergerfs at %s", mergerFs.Target)

	if err := unix.Unmount(mergerFs.Target, 0); err != nil {
		return errors.Wrap(err, "Failed to unmount mergerfs")
	}

	log.Debugf("MergerFs was successfully unmounted at %s", mergerFs.Target)

	return nil
}

func (mergerFs *MergerFsHandler) Cleanup() error {
	log.Tracef("Trying to cleanup mergerfs at %s", mergerFs.Target)

	if err := unix.Rmdir(mergerFs.Target); err != nil {
		return errors.Wrap(err, "Failed to remove mergerfs target dir")
	}

	log.Debugf("Removed mergerfs dir (%s)", mergerFs.Target)

	return nil
}
