package handler

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
)

type UnionFsHandler struct {
	Source       string
	MergeSources []string
	UpperDir     string
	Target       string
}

func NewUnionFs(source string, mergeSources []string, upperDir string, target string) *UnionFsHandler {
	return &UnionFsHandler{
		Source:       source,
		MergeSources: mergeSources,
		UpperDir:     upperDir,
		Target:       target,
	}
}

func (unionFs *UnionFsHandler) Mount() error {
	// TODO support multiple lower dirs
	log.Tracef("Trying to mount unionfs over %s into %s", unionFs.Source, unionFs.Target)

	if err := unionFs.setupFolders(); err != nil {
		return err
	}

	command := "/usr/bin/unionfs-fuse"
	args := []string{
		"-o",
		"cow,relaxed_permissions",
		fmt.Sprintf("%s=RW:%s=RO", unionFs.UpperDir, unionFs.Source),
		unionFs.Target,
	}

	if _, err := exec.Command(command, args...).CombinedOutput(); err != nil {
		return err
	}

	log.Debugf("mounted unionfs over %s into %s", unionFs.Source, unionFs.Target)

	return nil
}

func (unionFs *UnionFsHandler) setupFolders() error {
	if err := os.MkdirAll(unionFs.Target, 0755); err != nil {
		return errors.Wrap(err, "Failed to create unionFs target dir")
	}
	return nil
}

func (unionFs *UnionFsHandler) Unmount() error {
	log.Tracef("Trying to unmount unionfs at %s", unionFs.Target)

	// unix.Unmount(unionFs.Target, 0) is results in an "Operation not permitted" error
	if _, err := exec.Command("umount", unionFs.Target).CombinedOutput(); err != nil {
		return err
	}

	log.Debugf("UnionFs was successfully unmounted at %s", unionFs.Target)

	return nil
}

func (unionFs *UnionFsHandler) Cleanup() error {
	if err := unionFs.cleanupTargetDir(); err != nil {
		return err
	}
	return nil
}

func (unionFs *UnionFsHandler) cleanupTargetDir() error {
	log.Tracef("Trying to cleanup unionfs targert at %s", unionFs.Target)

	if err := unix.Rmdir(unionFs.Target); err != nil {
		fmt.Printf("Error removing unionfs target dir - %s\n", err)
		return err
	}
	log.Debugf("Removed unionfs target dir (%s)", unionFs.Target)
	return nil
}
