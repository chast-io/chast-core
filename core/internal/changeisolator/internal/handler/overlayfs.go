package handler

import (
	"fmt"
	"os"

	chastlog "chast.io/core/internal/logger"
	"github.com/joomcode/errorx"
	"golang.org/x/sys/unix"
)

type OverlayFsHandler struct {
	BaseDir  string
	Target   string
	UpperDir string
	WorkDir  string
}

func NewOverlayFs(baseDir string, target string, upperDir string, workDir string) *OverlayFsHandler {
	return &OverlayFsHandler{
		BaseDir:  baseDir,
		Target:   target,
		UpperDir: upperDir,
		WorkDir:  workDir,
	}
}

func (overlayFs *OverlayFsHandler) Mount() error {
	// TODO support multiple lower dirs
	chastlog.Log.Tracef("Trying to mount overlayfs over %s into %s", overlayFs.BaseDir, overlayFs.Target)

	if err := overlayFs.setupFolders(); err != nil {
		return err
	}

	fstype := "overlay"

	opts := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", overlayFs.BaseDir, overlayFs.UpperDir, overlayFs.WorkDir)
	chastlog.Log.Tracef("Mounting overlayfs with options: %s", opts)

	if err := unix.Mount("none", overlayFs.Target, fstype, unix.MS_NOSUID, opts); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to mount overlayfs")
	}

	chastlog.Log.Debugf("mounted overlayfs over %s into %s", overlayFs.BaseDir, overlayFs.Target)

	return nil
}

func (overlayFs *OverlayFsHandler) setupFolders() error {
	if _, err := os.Stat(overlayFs.BaseDir); os.IsNotExist(err) {
		return errorx.DataUnavailable.Wrap(err, "BaseDir does not exist")
	}

	if _, err := os.Stat(overlayFs.UpperDir); os.IsNotExist(err) {
		return errorx.DataUnavailable.Wrap(err, "UpperDir does not exist")
	}

	if err := os.MkdirAll(overlayFs.Target, 0o755); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to create overlayFs target dir")
	}

	if err := os.MkdirAll(overlayFs.WorkDir, 0o755); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to create overlayFs working dir")
	}

	return nil
}

func (overlayFs *OverlayFsHandler) Unmount() error {
	chastlog.Log.Tracef("Trying to unmount overlayfs at %s", overlayFs.Target)

	if err := unix.Unmount(overlayFs.Target, 0); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to unmount overlayfs")
	}

	chastlog.Log.Debugf("OverlayFs was successfully unmounted at %s", overlayFs.Target)

	return nil
}

func (overlayFs *OverlayFsHandler) Cleanup() error {
	if err := overlayFs.cleanupTargetDir(); err != nil {
		return err
	}

	if err := overlayFs.cleanupWorkingDir(); err != nil {
		return err
	}

	return nil
}

func (overlayFs *OverlayFsHandler) cleanupTargetDir() error {
	chastlog.Log.Tracef("Trying to cleanup overlayfs targert at %s", overlayFs.Target)

	if err := unix.Rmdir(overlayFs.Target); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to remove mergerfs dir")
	}

	chastlog.Log.Debugf("Removed overlayfs target dir (%s)", overlayFs.Target)

	return nil
}

func (overlayFs *OverlayFsHandler) cleanupWorkingDir() error {
	chastlog.Log.Tracef("Trying to cleanup overlayfs working dir at %s", overlayFs.WorkDir)

	if err := os.RemoveAll(overlayFs.WorkDir); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to remove overlayfs working dir")
	}

	chastlog.Log.Debugf("Removed overlayfs working dir (%s)", overlayFs.WorkDir)

	return nil
}
