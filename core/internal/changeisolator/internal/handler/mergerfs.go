package handler

import (
	"os"
	"os/exec"

	chastlog "chast.io/core/internal/logger"
	"github.com/joomcode/errorx"
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
	chastlog.Log.Tracef("Trying to merge %s into %s", mergerFs.Source, mergerFs.Target)

	if err := os.MkdirAll(mergerFs.Target, 0o755); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to create mergerFs target dir")
	}

	command := "/usr/bin/mergerfs"
	args := []string{
		mergerFs.Source,
		mergerFs.Target,
	}

	if _, err := exec.Command(command, args...).CombinedOutput(); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to mount mergerfs")
	}

	chastlog.Log.Debugf("Mergerfs was successfully mounted -  %s into %s", mergerFs.Source, mergerFs.Target)

	return nil
}

func (mergerFs *MergerFsHandler) Unmount() error {
	chastlog.Log.Tracef("Trying to unmerge mergerfs at %s", mergerFs.Target)

	if err := unix.Unmount(mergerFs.Target, 0); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to unmount mergerfs")
	}

	chastlog.Log.Debugf("MergerFs was successfully unmounted at %s", mergerFs.Target)

	return nil
}

func (mergerFs *MergerFsHandler) Cleanup() error {
	chastlog.Log.Tracef("Trying to cleanup mergerfs at %s", mergerFs.Target)

	if err := unix.Rmdir(mergerFs.Target); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to remove mergerfs target dir")
	}

	chastlog.Log.Debugf("Removed mergerfs dir (%s)", mergerFs.Target)

	return nil
}
