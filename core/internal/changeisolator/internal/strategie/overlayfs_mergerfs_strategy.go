package strategie

import (
	"path/filepath"

	"chast.io/core/internal/changeisolator/internal/handler"
	"chast.io/core/internal/changeisolator/pkg/strategy"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type OverlayFsMergerFsStrategy struct {
	IsolatorContext

	mergerFsHandler   *handler.MergerFsHandler
	overlayFsHandler  *handler.OverlayFsHandler
	devMounter        *handler.Mounter
	procMounter       *handler.Mounter
	tmpMounter        *handler.TmpMounter
	changeRootHandler *handler.ChangeRootHandler
}

func NewOverlayFsMergerFsStrategy(context IsolatorContext) *OverlayFsMergerFsStrategy {
	return &OverlayFsMergerFsStrategy{ //nolint:exhaustruct // handler initialization is done in Initialize method
		IsolatorContext: context,
	}
}

func (strat *OverlayFsMergerFsStrategy) GetIsolationStrategy() strategy.IsolationStrategy {
	return strategy.OverlayFS
}

func (strat *OverlayFsMergerFsStrategy) Initialize() error {
	log.Tracef("Initializing change isolator with the overlayfs mergerfs strategy")

	if err := strat.IsolatorContext.Initialize(); err != nil {
		return err
	}

	rootFolder := strat.RootFolder
	newRootFsFolder := filepath.Join(strat.OperationDirectory, "rootfs")
	mergerFsFolder := filepath.Join(strat.OperationDirectory, "mergerfs")
	overlayFsWorkingDirFolder := filepath.Join(strat.OperationDirectory, "overlayFsWorkingDir")

	strat.mergerFsHandler = handler.NewMergerFs(rootFolder, mergerFsFolder)

	strat.overlayFsHandler = handler.NewOverlayFs(
		mergerFsFolder,
		newRootFsFolder,
		strat.ChangeCaptureFolder,
		overlayFsWorkingDirFolder,
	)
	strat.devMounter = handler.NewMounter("dev", rootFolder, newRootFsFolder)
	strat.procMounter = handler.NewMounter("proc", rootFolder, newRootFsFolder)
	strat.tmpMounter = handler.NewTmpMounter(rootFolder, newRootFsFolder, strat.OperationDirectory)

	strat.changeRootHandler = handler.NewChangeRoot(newRootFsFolder, strat.WorkingDirectory)

	return nil
}

// === Prepare ===

func (strat *OverlayFsMergerFsStrategy) PrepareOutsideNS() error {
	log.Tracef("[Outside NS] Preparing change isolator with the overlayfs mergerfs strategy")

	return nil
}

func (strat *OverlayFsMergerFsStrategy) PrepareInsideNS() error {
	log.Tracef("[Inside NS] Preparing change isolator with the overlayfs mergerfs strategy")

	if err := strat.mergerFsHandler.Mount(); err != nil {
		return errors.Wrap(err, "Error mounting mergerfs")
	}

	if err := strat.overlayFsHandler.Mount(); err != nil {
		return errors.Wrap(err, "Error mounting overlayfs")
	}

	if err := strat.devMounter.Mount(); err != nil {
		return errors.Wrap(err, "Error mounting dev")
	}

	if err := strat.procMounter.Mount(); err != nil {
		return errors.Wrap(err, "Error mounting proc")
	}

	if err := strat.tmpMounter.Mount(); err != nil {
		return errors.Wrap(err, "Error mounting tmp")
	}

	if err := strat.changeRootHandler.Init(); err != nil {
		return errors.Wrap(err, "Error initializing change root")
	}

	if err := strat.changeRootHandler.Open(); err != nil {
		return errors.Wrap(err, "Error open change root")
	}

	return nil
}

// === Cleanup ===

func (strat *OverlayFsMergerFsStrategy) CleanupInsideNS() error {
	log.Tracef("[Inside NS] Cleaning up change isolator with the overlayfs mergerfs strategy")

	if err := strat.changeRootHandler.Close(); err != nil {
		return errors.Wrap(err, "Error closing change root")
	}

	if err := strat.tmpMounter.Unmount(); err != nil {
		return errors.Wrap(err, "Error unmounting tmp")
	}

	if err := strat.procMounter.Unmount(); err != nil {
		return errors.Wrap(err, "Error unmounting proc")
	}

	if err := strat.devMounter.Unmount(); err != nil {
		return errors.Wrap(err, "Error unmounting dev")
	}

	if err := strat.overlayFsHandler.Unmount(); err != nil {
		return errors.Wrap(err, "Error unmounting overlayfs")
	}

	if err := strat.overlayFsHandler.Cleanup(); err != nil {
		return errors.Wrap(err, "Error cleaning up overlayfs")
	}

	if err := strat.mergerFsHandler.Unmount(); err != nil {
		return errors.Wrap(err, "Error unmounting mergerfs")
	}

	if err := strat.mergerFsHandler.Cleanup(); err != nil {
		return errors.Wrap(err, "Error cleaning up mergerfs")
	}

	return strat.IsolatorContext.CleanupInsideNS()
}

func (strat *OverlayFsMergerFsStrategy) CleanupOutsideNS() error {
	log.Tracef("[Outside NS] Cleaning up change isolator with the overlayfs mergerfs strategy")

	return strat.IsolatorContext.CleanupOutsideNS()
}
