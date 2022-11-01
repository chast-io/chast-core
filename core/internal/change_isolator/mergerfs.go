package change_isolator

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
)

type mergerFsHandler struct {
	Source string
	Target string
}

func newMergerFs(source string, target string) *mergerFsHandler {
	return &mergerFsHandler{
		Source: source,
		Target: target,
	}
}

func (mergerFs *mergerFsHandler) mount() error {
	log.Tracef("Trying to merge %s into %s", mergerFs.Source, mergerFs.Target)

	if err := os.MkdirAll(mergerFs.Target, 0755); err != nil {
		return errors.Wrap(err, "Failed to create mergerFs target dir")
	}

	command := "/usr/bin/mergerfs"
	args := []string{
		mergerFs.Source,
		mergerFs.Target,
	}

	if _, err := exec.Command(command, args...).CombinedOutput(); err != nil {
		return err
	}

	log.Debugf("Mergerfs was successfully mounted -  %s into %s", mergerFs.Source, mergerFs.Target)

	return nil
}

func (mergerFs *mergerFsHandler) unmount() error {
	log.Tracef("Trying to unmerge mergerfs at %s", mergerFs.Target)

	if err := unix.Unmount(mergerFs.Target, 0); err != nil {
		return err
	}

	log.Debugf("MergerFs was successfully unmounted at %s", mergerFs.Target)

	return nil
}

func (mergerFs *mergerFsHandler) cleanup() error {
	log.Tracef("Trying to cleanup mergerfs at %s", mergerFs.Target)

	if err := unix.Rmdir(mergerFs.Target); err != nil {
		return err
	}
	log.Debugf("Removed mergerfs dir (%s)", mergerFs.Target)
	return nil
}
