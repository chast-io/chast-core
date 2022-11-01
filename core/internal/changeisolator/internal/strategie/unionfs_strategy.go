package strategie

import (
	"chast.io/core/internal/changeisolator/internal/handler"
	"chast.io/core/internal/changeisolator/pkg/strategy"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

type UnionFsStrategy struct {
	IsolatorContext

	unionFsHandler    *handler.UnionFsHandler
	devMounter        *handler.Mounter
	procMounter       *handler.Mounter
	changeRootHandler *handler.ChangeRootHandler
}

func NewUnionFsStrategy(context IsolatorContext) *UnionFsStrategy {
	return &UnionFsStrategy{
		IsolatorContext: context,
	}
}

func (strat *UnionFsStrategy) GetIsolationStrategy() strategy.IsolationStrategy {
	return strategy.UnionFS
}

// === Prepare ===

func (strat *UnionFsStrategy) Initialize() error {
	log.Tracef("Initializing change isolator with the unionfs strategy")

	if err := strat.IsolatorContext.Initialize(); err != nil {
		return err
	}

	rootFolder := strat.RootFolder
	newRootFsFolder := filepath.Join(strat.OperationDirectory, "rootfs")

	strat.unionFsHandler = handler.NewUnionFs(rootFolder, []string{}, strat.ChangeCaptureFolder, newRootFsFolder)

	strat.devMounter = handler.NewMounter("dev", rootFolder, newRootFsFolder)
	strat.procMounter = handler.NewMounter("proc", rootFolder, newRootFsFolder)

	strat.changeRootHandler = handler.NewChangeRoot(newRootFsFolder, strat.WorkingDirectory)

	return nil
}

func (strat *UnionFsStrategy) PrepareOutsideNS() error {
	log.Tracef("[Outside NS] Preparing change isolator with the unionfs strategy")

	if err := strat.unionFsHandler.Mount(); err != nil {
		return errors.Wrap(err, "Error mounting unionfs")
	}

	return nil
}

func (strat *UnionFsStrategy) PrepareInsideNS() error {
	log.Tracef("[Inside NS] Preparing change isolator with the unionfs strategy")

	// TODO mount empty /tmp folder to prevent recursive alterations and to provide a clean and temporary tmp folder

	if err := strat.devMounter.Mount(); err != nil {
		return errors.Wrap(err, "Error mounting dev")
	}

	if err := strat.procMounter.Mount(); err != nil {
		return errors.Wrap(err, "Error mounting proc")
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

func (strat *UnionFsStrategy) CleanupInsideNS() error {
	log.Tracef("[Inside NS] Cleaning up change isolator with the unionfs strategy")

	if err := strat.changeRootHandler.Close(); err != nil {
		return errors.Wrap(err, "Error closing change root")
	}

	if err := strat.procMounter.Unmount(); err != nil {
		return errors.Wrap(err, "Error unmounting proc")
	}

	if err := strat.devMounter.Unmount(); err != nil {
		return errors.Wrap(err, "Error unmounting dev")
	}

	return strat.IsolatorContext.CleanupInsideNS()
}

func (strat *UnionFsStrategy) CleanupOutsideNS() error {
	log.Tracef("[Outside NS] Cleaning up change isolator with the unionfs strategy")

	if err := strat.unionFsHandler.Unmount(); err != nil {
		return errors.Wrap(err, "Error unmounting unionfs")
	}

	if err := strat.unionFsHandler.Cleanup(); err != nil {
		return errors.Wrap(err, "Error cleaning up unionfs")
	}

	return strat.IsolatorContext.CleanupOutsideNS()
}
