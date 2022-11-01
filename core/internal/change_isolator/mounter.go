package change_isolator

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"path/filepath"
	"syscall"
)

type mounter struct {
	Folder string
	Source string
	Target string
}

func newMounter(folder string, root string, target string) *mounter {
	return &mounter{
		Folder: folder,
		Source: filepath.Join(root, folder),
		Target: filepath.Join(target, folder),
	}
}

func (mntr *mounter) mount() error {
	log.Tracef("Trying to mount %s into %s", mntr.Source, mntr.Target)
	fstype := ""
	if mntr.Folder == "proc" {
		fstype = "proc"
	}
	flags := unix.MS_BIND | unix.MS_REC | unix.MS_PRIVATE | unix.MS_RDONLY
	data := ""

	if err := syscall.Mount(mntr.Source, mntr.Target, fstype, uintptr(flags), data); err != nil {
		return err
	}

	log.Tracef("Mounted %s into %s", mntr.Source, mntr.Target)
	return nil
}

func (mntr *mounter) unmount() error {
	log.Tracef("Trying to unmount %s lazily", mntr.Target)
	if errLazy := syscall.Unmount(mntr.Target, unix.MNT_DETACH); errLazy != nil {
		return errLazy
	}

	log.Tracef("Unmounte %s", mntr.Target)

	return nil
}
