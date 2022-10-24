package main

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
	"syscall"
)

type mergerFsHandle struct {
	Source string
	Target string
}

func newMergerFs(source string, target string) *mergerFsHandle {
	return &mergerFsHandle{
		Source: source,
		Target: target,
	}
}

func (mergerFs *mergerFsHandle) mount() error {
	log.Tracef("Trying to merge %s into %s", mergerFs.Source, mergerFs.Target)

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

func (mergerFs *mergerFsHandle) unmount() error {
	log.Tracef("Trying to unmerge mergerfs at %s", mergerFs.Target)

	if err := syscall.Unmount(mergerFs.Target, 0); err != nil {
		return err
	}

	log.Debugf("MergerFs was successfully unmounted at %s", mergerFs.Target)

	return nil
}

func (mergerFs *mergerFsHandle) cleanup() error {
	log.Tracef("Trying to cleanup mergerfs at %s", mergerFs.Target)

	if err := syscall.Rmdir(mergerFs.Target); err != nil {
		return err
	}
	log.Debugf("Removed mergerfs dir (%s)", mergerFs.Target)
	return nil
}
