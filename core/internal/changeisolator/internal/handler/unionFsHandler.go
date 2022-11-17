package handler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	chastlog "chast.io/core/internal/logger"
	"github.com/joomcode/errorx"
	"golang.org/x/sys/unix"
)

type UnionFsHandler struct {
	Source            string
	SourceJoinFolders []string
	UpperDir          string
	Target            string
}

func NewUnionFs(source string, sourceJoinFolders []string, upperDir string, target string) *UnionFsHandler {
	return &UnionFsHandler{
		Source:            source,
		SourceJoinFolders: sourceJoinFolders,
		UpperDir:          upperDir,
		Target:            target,
	}
}

func (unionFs *UnionFsHandler) Mount() error {
	chastlog.Log.Tracef("Trying to mount unionfs over %s into %s", unionFs.Source, unionFs.Target)

	if err := unionFs.setupFolders(); err != nil {
		return err
	}

	command := "/usr/bin/unionfs-fuse"

	args := []string{
		"-o",
		// cow             : copy on write, changes are only reflected in the upper dir
		// hide_meta_files : hide .unionfs-fuse folder present in each dir except the lower dir
		//				     and containing all the metadata (removal information)
		// hard_remove     : remove files from the upper dir instead of just hiding them
		"cow,hide_meta_files,hard_remove",
		fmt.Sprintf("%s=RW:%s", unionFs.UpperDir, unionFs.getSourceFolders()),
		unionFs.Target,
	}

	if _, err := exec.Command(command, args...).CombinedOutput(); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to mount unionfs")
	}

	chastlog.Log.Debugf("mounted unionfs over %s into %s", unionFs.Source, unionFs.Target)

	return nil
}

func (unionFs *UnionFsHandler) getSourceFolders() string {
	var sourceFolders strings.Builder

	size := len(unionFs.SourceJoinFolders)

	for i := range unionFs.SourceJoinFolders {
		// reverse order required to have the last task be the left most and the first task be the right most path
		sourceFolders.WriteString(fmt.Sprintf("%s=RO:", unionFs.SourceJoinFolders[size-1-i]))
	}

	sourceFolders.WriteString(fmt.Sprintf("%s=RO", unionFs.Source))

	return sourceFolders.String()
}

func (unionFs *UnionFsHandler) setupFolders() error {
	if err := os.MkdirAll(unionFs.Target, 0o755); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to create unionFs target dir")
	}

	return nil
}

func (unionFs *UnionFsHandler) Unmount() error {
	chastlog.Log.Tracef("Trying to unmount unionfs at %s", unionFs.Target)

	// unix.Unmount(unionFs.Target, 0) is results in an "Operation not permitted" error
	if _, err := exec.Command("umount", unionFs.Target).CombinedOutput(); err != nil { //nolint:gosec // secure
		return errorx.ExternalError.Wrap(err, "Failed to unmount unionfs")
	}

	chastlog.Log.Debugf("UnionFs was successfully unmounted at %s", unionFs.Target)

	return nil
}

func (unionFs *UnionFsHandler) Cleanup() error {
	if err := unionFs.cleanupTargetDir(); err != nil {
		return err
	}

	return nil
}

func (unionFs *UnionFsHandler) cleanupTargetDir() error {
	chastlog.Log.Tracef("Trying to cleanup unionfs targert at %s", unionFs.Target)

	if err := unix.Rmdir(unionFs.Target); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to remove unionfs target dir")
	}

	chastlog.Log.Debugf("Removed unionfs target dir (%s)", unionFs.Target)

	return nil
}
