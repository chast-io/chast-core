package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"syscall"
)

type overlayFsHandle struct {
	RootFsPath string
	Target     string
	UpperDir   string
	WorkDir    string
}

func newOverlayFs(rootFsDir string, target string, upperDir string, workDir string) *overlayFsHandle {
	return &overlayFsHandle{
		RootFsPath: rootFsDir,
		Target:     target,
		UpperDir:   upperDir,
		WorkDir:    workDir,
	}
}

func (overlayFs *overlayFsHandle) mount() error {
	log.Tracef("trying to mount overlayfs over %s into %s", overlayFs.RootFsPath, overlayFs.Target)

	fstype := "overlay"

	opts := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", overlayFs.RootFsPath, overlayFs.UpperDir, overlayFs.WorkDir)
	if err := syscall.Mount("none", overlayFs.Target, fstype, syscall.MS_NOSUID, opts); err != nil {
		return err
	}

	log.Debugf("mounted overlayfs over %s into %s", overlayFs.RootFsPath, overlayFs.Target)

	return nil
}

func (overlayFs *overlayFsHandle) unmount() error {
	log.Tracef("Trying to unmerge overlayfs at %s", overlayFs.Target)

	if err := syscall.Unmount(overlayFs.Target, 0); err != nil {
		return err
	}

	log.Debugf("OverlayFs was successfully unmounted at %s", overlayFs.Target)

	return nil
}

func (overlayFs *overlayFsHandle) cleanup() error {
	if err := overlayFs.cleanupTargetDir(); err != nil {
		return err
	}
	if err := overlayFs.cleanupWorkingDir(); err != nil {
		return err
	}
	return nil
}

func (overlayFs *overlayFsHandle) cleanupTargetDir() error {
	log.Tracef("Trying to cleanup overlayfs targert at %s", overlayFs.Target)

	if err := syscall.Rmdir(overlayFs.Target); err != nil {
		fmt.Printf("Error removing mergerfs dir - %s\n", err)
		return err
	}
	log.Debugf("Removed overlayfs target dir (%s)", overlayFs.Target)
	return nil
}

func (overlayFs *overlayFsHandle) cleanupWorkingDir() error {
	log.Tracef("Trying to cleanup overlayfs working dir at %s", overlayFs.WorkDir)

	if err := os.RemoveAll(overlayFs.WorkDir); err != nil {
		return err
	}

	log.Debugf("Removed overlayfs working dir (%s)", overlayFs.WorkDir)
	return nil
}
