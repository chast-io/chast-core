package overlay

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

type changeIsolatorOverlayfsMergerfsStrategy struct {
	changeIsolator

	mergerFsHandler   *mergerFsHandler
	overlayFsHandler  *overlayFsHandler
	devMounter        *mounter
	procMounter       *mounter
	changeRootHandler *changeRootHandler
}

func newChangeIsolatorOverlayfsMergerfsStrategy(changeIsolator changeIsolator) *changeIsolatorOverlayfsMergerfsStrategy {
	return &changeIsolatorOverlayfsMergerfsStrategy{
		changeIsolator: changeIsolator,
	}
}

func (strategy *changeIsolatorOverlayfsMergerfsStrategy) getIsolationStrategy() IsolationStrategy {
	return OverlayFS
}

func (strategy *changeIsolatorOverlayfsMergerfsStrategy) initialize() error {
	log.Tracef("Initializing change isolator with the overlayfs mergerfs strategy")

	if err := strategy.changeIsolator.initialize(); err != nil {
		return err
	}

	rootFolder := strategy.RootFolder
	newRootFsFolder := filepath.Join(strategy.OperationDirectory, "rootfs")
	mergerFsFolder := filepath.Join(strategy.OperationDirectory, "mergerfs")
	overlayFsWorkingDirFolder := filepath.Join(strategy.OperationDirectory, "overlayFsWorkingDir")

	strategy.mergerFsHandler = newMergerFs(rootFolder, mergerFsFolder)

	strategy.overlayFsHandler = newOverlayFs(mergerFsFolder, newRootFsFolder, strategy.ChangeCaptureFolder, overlayFsWorkingDirFolder)
	strategy.devMounter = newMounter("dev", rootFolder, newRootFsFolder)
	strategy.procMounter = newMounter("proc", rootFolder, newRootFsFolder)

	strategy.changeRootHandler = newChangeRoot(newRootFsFolder, strategy.WorkingDirectory)

	return nil
}

// === Prepare ===

func (strategy *changeIsolatorOverlayfsMergerfsStrategy) prepareOutsideNS() error {
	log.Tracef("[Outside NS] Preparing change isolator with the overlayfs mergerfs strategy")
	return nil

}

func (strategy *changeIsolatorOverlayfsMergerfsStrategy) prepareInsideNS() error {
	log.Tracef("[Inside NS] Preparing change isolator with the overlayfs mergerfs strategy")

	if err := strategy.mergerFsHandler.mount(); err != nil {
		return errors.Wrap(err, "Error mounting mergerfs")
	}

	if err := strategy.overlayFsHandler.mount(); err != nil {
		return errors.Wrap(err, "Error mounting overlayfs")
	}

	// TODO mount empty /tmp folder to prevent recursive alterations and to provide a clean and temporary tmp folder

	if err := strategy.devMounter.mount(); err != nil {
		return errors.Wrap(err, "Error mounting dev")
	}

	if err := strategy.procMounter.mount(); err != nil {
		return errors.Wrap(err, "Error mounting proc")
	}

	if err := strategy.changeRootHandler.init(); err != nil {
		return errors.Wrap(err, "Error initializing change root")
	}

	if err := strategy.changeRootHandler.open(); err != nil {
		return errors.Wrap(err, "Error open change root")
	}

	return nil
}

// === Cleanup ===

func (strategy *changeIsolatorOverlayfsMergerfsStrategy) cleanupInsideNS() error {
	log.Tracef("[Inside NS] Cleaning up change isolator with the overlayfs mergerfs strategy")

	if err := strategy.changeRootHandler.close(); err != nil {
		return errors.Wrap(err, "Error closing change root")
	}

	if err := strategy.procMounter.unmount(); err != nil {
		return errors.Wrap(err, "Error unmounting proc")
	}

	if err := strategy.devMounter.unmount(); err != nil {
		return errors.Wrap(err, "Error unmounting dev")
	}

	if err := strategy.overlayFsHandler.unmount(); err != nil {
		return errors.Wrap(err, "Error unmounting overlayfs")
	}

	if err := strategy.overlayFsHandler.cleanup(); err != nil {
		return errors.Wrap(err, "Error cleaning up overlayfs")
	}

	if err := strategy.mergerFsHandler.unmount(); err != nil {
		return errors.Wrap(err, "Error unmounting mergerfs")
	}

	if err := strategy.mergerFsHandler.cleanup(); err != nil {
		return errors.Wrap(err, "Error cleaning up mergerfs")
	}

	return strategy.changeIsolator.cleanupInsideNS()
}

func (strategy *changeIsolatorOverlayfsMergerfsStrategy) cleanupOutsideNS() error {
	log.Tracef("[Outside NS] Cleaning up change isolator with the overlayfs mergerfs strategy")

	return strategy.changeIsolator.cleanupOutsideNS()
}
