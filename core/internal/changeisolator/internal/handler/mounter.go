package handler

import (
	"path/filepath"
	"syscall"

	chastlog "chast.io/core/internal/logger"
	"github.com/joomcode/errorx"
	"golang.org/x/sys/unix"
)

type Mounter struct {
	Folder string
	Source string
	Target string
}

func NewMounter(folder string, root string, target string) *Mounter {
	return &Mounter{
		Folder: folder,
		Source: filepath.Join(root, folder),
		Target: filepath.Join(target, folder),
	}
}

func (mntr *Mounter) Mount() error {
	chastlog.Log.Tracef("Trying to mount %s into %s", mntr.Source, mntr.Target)

	fstype := ""

	if mntr.Folder == "proc" {
		fstype = "proc"
	}

	flags := unix.MS_BIND | unix.MS_REC | unix.MS_PRIVATE | unix.MS_RDONLY
	data := ""

	if err := syscall.Mount(mntr.Source, mntr.Target, fstype, uintptr(flags), data); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to mount")
	}

	chastlog.Log.Tracef("Mounted %s into %s", mntr.Source, mntr.Target)

	return nil
}

func (mntr *Mounter) Unmount() error {
	chastlog.Log.Tracef("Trying to unmount %s lazily", mntr.Target)

	if err := syscall.Unmount(mntr.Target, unix.MNT_DETACH); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to unmount lazily")
	}

	chastlog.Log.Tracef("Unmounte %s", mntr.Target)

	return nil
}
