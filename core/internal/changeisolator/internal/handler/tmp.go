package handler

import (
	"os"
	"path/filepath"

	"github.com/joomcode/errorx"
)

type TmpMounter struct {
	mounter Mounter

	OperationDirectory string
}

func NewTmpMounter(source string, target string, operationDirectory string) *TmpMounter {
	return &TmpMounter{
		mounter:            *NewMounter("tmp", source, target),
		OperationDirectory: operationDirectory,
	}
}

func (mntr *TmpMounter) Mount() error {
	if err := os.Mkdir(filepath.Join(mntr.OperationDirectory, "tmp"), 0755); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to create tmp directory")
	}

	return mntr.mounter.Mount()
}

func (mntr *TmpMounter) Unmount() error {
	if err := os.RemoveAll(filepath.Join(mntr.OperationDirectory, "tmp")); err != nil {
		return errorx.ExternalError.Wrap(err, "Failed to remove tmp directory")
	}

	return mntr.mounter.Unmount()
}
