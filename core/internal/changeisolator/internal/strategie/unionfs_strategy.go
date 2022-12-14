package strategie

import (
	"path/filepath"

	"chast.io/core/internal/changeisolator/internal/handler"
	"chast.io/core/internal/changeisolator/pkg/strategy"
	chastlog "chast.io/core/internal/logger"
	"github.com/joomcode/errorx"
)

type UnionFsStrategy struct {
	IsolatorContext

	unionFsHandler    *handler.UnionFsHandler
	devMounter        *handler.Mounter
	procMounter       *handler.Mounter
	tmpMounter        *handler.TmpMounter
	changeRootHandler *handler.ChangeRootHandler
}

func NewUnionFsStrategy(context IsolatorContext) *UnionFsStrategy {
	return &UnionFsStrategy{ //nolint:exhaustruct // handlers are initialized separately and explicitly
		IsolatorContext: context,
	}
}

func (strat *UnionFsStrategy) GetIsolationStrategy() strategy.IsolationStrategy {
	return strategy.UnionFS
}

// === Prepare ===

func (strat *UnionFsStrategy) Initialize() error {
	chastlog.Log.Tracef("Initializing change isolator with the unionfs strategy")

	if err := strat.IsolatorContext.Initialize(); err != nil {
		return err
	}

	rootFolder := strat.RootFolder
	newRootFsFolder := filepath.Join(strat.OperationDirectory, "rootfs")

	strat.unionFsHandler = handler.NewUnionFs(
		rootFolder,
		strat.RootJoinFolders,
		strat.ChangeCaptureFolder,
		newRootFsFolder,
	)

	strat.devMounter = handler.NewMounter("dev", rootFolder, newRootFsFolder)
	strat.procMounter = handler.NewMounter("proc", rootFolder, newRootFsFolder)
	strat.tmpMounter = handler.NewTmpMounter(rootFolder, newRootFsFolder, strat.OperationDirectory)

	strat.changeRootHandler = handler.NewChangeRoot(newRootFsFolder, strat.WorkingDirectory)

	return nil
}

func (strat *UnionFsStrategy) PrepareOutsideNS() error {
	chastlog.Log.Tracef("[Outside NS] Preparing change isolator with the unionfs strategy")

	if err := strat.unionFsHandler.Mount(); err != nil {
		return errorx.InternalError.Wrap(err, "Error mounting unionfs")
	}

	return nil
}

func (strat *UnionFsStrategy) PrepareInsideNS() error {
	chastlog.Log.Tracef("[Inside NS] Preparing change isolator with the unionfs strategy")

	if err := strat.devMounter.Mount(); err != nil {
		return errorx.InternalError.Wrap(err, "Error mounting dev")
	}

	if err := strat.procMounter.Mount(); err != nil {
		return errorx.InternalError.Wrap(err, "Error mounting proc")
	}

	if err := strat.tmpMounter.Mount(); err != nil {
		return errorx.InternalError.Wrap(err, "Error mounting tmp")
	}

	if err := strat.changeRootHandler.Init(); err != nil {
		return errorx.InternalError.Wrap(err, "Error initializing change root")
	}

	if err := strat.changeRootHandler.Open(); err != nil {
		return errorx.InternalError.Wrap(err, "Error open change root")
	}

	return nil
}

// === Cleanup ===

func (strat *UnionFsStrategy) CleanupInsideNS() error {
	chastlog.Log.Tracef("[Inside NS] Cleaning up change isolator with the unionfs strategy")

	if err := strat.changeRootHandler.Close(); err != nil {
		return errorx.InternalError.Wrap(err, "Error closing change root")
	}

	if err := strat.tmpMounter.Unmount(); err != nil {
		return errorx.InternalError.Wrap(err, "Error unmounting tmp")
	}

	if err := strat.procMounter.Unmount(); err != nil {
		return errorx.InternalError.Wrap(err, "Error unmounting proc")
	}

	if err := strat.devMounter.Unmount(); err != nil {
		return errorx.InternalError.Wrap(err, "Error unmounting dev")
	}

	return strat.IsolatorContext.CleanupInsideNS()
}

func (strat *UnionFsStrategy) CleanupOutsideNS() error {
	chastlog.Log.Tracef("[Outside NS] Cleaning up change isolator with the unionfs strategy")

	if err := strat.unionFsHandler.Unmount(); err != nil {
		return errorx.InternalError.Wrap(err, "Error unmounting unionfs")
	}

	if err := strat.unionFsHandler.Cleanup(); err != nil {
		return errorx.InternalError.Wrap(err, "Error cleaning up unionfs")
	}

	return strat.IsolatorContext.CleanupOutsideNS()
}
