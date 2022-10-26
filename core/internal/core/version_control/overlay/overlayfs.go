package overlay

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"os"
)

type overlayFsHandler struct {
	BaseDir  string
	Target   string
	UpperDir string
	WorkDir  string
}

func newOverlayFs(baseDir string, target string, upperDir string, workDir string) *overlayFsHandler {
	return &overlayFsHandler{
		BaseDir:  baseDir,
		Target:   target,
		UpperDir: upperDir,
		WorkDir:  workDir,
	}
}

func (overlayFs *overlayFsHandler) setupFolders() error {
	if _, err := os.Stat(overlayFs.BaseDir); os.IsNotExist(err) {
		return errors.Wrap(err, "BaseDir does not exist")
	}

	if _, err := os.Stat(overlayFs.UpperDir); os.IsNotExist(err) {
		return errors.Wrap(err, "UpperDir does not exist")
	}

	if err := os.MkdirAll(overlayFs.Target, 0755); err != nil {
		return errors.Wrap(err, "Failed to create overlayFs target dir")
	}

	if err := os.MkdirAll(overlayFs.WorkDir, 0755); err != nil {
		return errors.Wrap(err, "Failed to create overlayFs working dir")
	}

	return nil
}

func (overlayFs *overlayFsHandler) mount() error {
	log.Tracef("Trying to mount overlayfs over %s into %s", overlayFs.BaseDir, overlayFs.Target)

	if err := overlayFs.setupFolders(); err != nil {
		return err
	}

	fstype := "overlay"

	opts := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", overlayFs.BaseDir, overlayFs.UpperDir, overlayFs.WorkDir)
	if err := unix.Mount("none", overlayFs.Target, fstype, unix.MS_NOSUID, opts); err != nil {
		return err
	}

	log.Debugf("mounted overlayfs over %s into %s", overlayFs.BaseDir, overlayFs.Target)

	return nil
}

func (overlayFs *overlayFsHandler) unmount() error {
	log.Tracef("Trying to unmount overlayfs at %s", overlayFs.Target)

	if err := unix.Unmount(overlayFs.Target, 0); err != nil {
		return err
	}

	log.Debugf("OverlayFs was successfully unmounted at %s", overlayFs.Target)

	return nil
}

func (overlayFs *overlayFsHandler) cleanup() error {
	if err := overlayFs.cleanupTargetDir(); err != nil {
		return err
	}
	if err := overlayFs.cleanupWorkingDir(); err != nil {
		return err
	}
	return nil
}

func (overlayFs *overlayFsHandler) cleanupTargetDir() error {
	log.Tracef("Trying to cleanup overlayfs targert at %s", overlayFs.Target)

	if err := unix.Rmdir(overlayFs.Target); err != nil {
		fmt.Printf("Error removing mergerfs dir - %s\n", err)
		return err
	}
	log.Debugf("Removed overlayfs target dir (%s)", overlayFs.Target)
	return nil
}

func (overlayFs *overlayFsHandler) cleanupWorkingDir() error {
	log.Tracef("Trying to cleanup overlayfs working dir at %s", overlayFs.WorkDir)

	if err := os.RemoveAll(overlayFs.WorkDir); err != nil {
		return err
	}

	log.Debugf("Removed overlayfs working dir (%s)", overlayFs.WorkDir)
	return nil
}
