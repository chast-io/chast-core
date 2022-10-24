package main

import (
	"os"
	"syscall"
)

type changeRootHandler struct {
	RootFsPath                 string
	originalRootFileDescriptor *os.File
}

func newChangeRoot(rootFsPath string) *changeRootHandler {
	return &changeRootHandler{
		RootFsPath: rootFsPath,
	}
}

func (changeRoot *changeRootHandler) init() error {
	root, err := os.Open("/")
	if err != nil {
		return err
	}
	changeRoot.originalRootFileDescriptor = root

	return nil
}

func (changeRoot *changeRootHandler) open() error {
	if err := syscall.Chroot(changeRoot.RootFsPath); err != nil {
		if err := changeRoot.close(); err != nil {
			return err
		}
		return err
	}
	if err := syscall.Chdir("/"); err != nil {
		return err
	}

	return nil
}

func (changeRoot *changeRootHandler) close() error {
	err := changeRoot.originalRootFileDescriptor.Close()
	if err != nil {
		return err
	}
	return nil
}
