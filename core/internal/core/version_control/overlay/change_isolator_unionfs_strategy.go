package overlay

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

type changeIsolatorUnionFsStrategy struct {
	changeIsolator

	unionFsHandler    *unionFsHandler
	devMounter        *mounter
	procMounter       *mounter
	changeRootHandler *changeRootHandler
}

func newChangeIsolatorUnionFsStrategy(changeIsolator changeIsolator) *changeIsolatorUnionFsStrategy {
	return &changeIsolatorUnionFsStrategy{
		changeIsolator: changeIsolator,
	}
}

func (strategy *changeIsolatorUnionFsStrategy) getIsolationStrategy() IsolationStrategy {
	return UnionFS
}

// === Prepare ===

func (strategy *changeIsolatorUnionFsStrategy) initialize() error {
	log.Tracef("Initializing change isolator with the unionfs strategy")

	if err := strategy.changeIsolator.initialize(); err != nil {
		return err
	}

	rootFolder := strategy.RootFolder
	newRootFsFolder := filepath.Join(strategy.OperationDirectory, "rootfs")

	strategy.unionFsHandler = newUnionFs(rootFolder, []string{}, strategy.ChangeCaptureFolder, newRootFsFolder)

	strategy.devMounter = newMounter("dev", rootFolder, newRootFsFolder)
	strategy.procMounter = newMounter("proc", rootFolder, newRootFsFolder)

	strategy.changeRootHandler = newChangeRoot(newRootFsFolder, strategy.WorkingDirectory)

	return nil
}

func (strategy *changeIsolatorUnionFsStrategy) prepareOutsideNS() error {
	log.Tracef("[Outside NS] Preparing change isolator with the unionfs strategy")

	if err := strategy.unionFsHandler.mount(); err != nil {
		return errors.Wrap(err, "Error mounting unionfs")
	}

	return nil
}

func (strategy *changeIsolatorUnionFsStrategy) prepareInsideNS() error {
	log.Tracef("[Inside NS] Preparing change isolator with the unionfs strategy")

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

func (strategy *changeIsolatorUnionFsStrategy) cleanupInsideNS() error {
	log.Tracef("[Inside NS] Cleaning up change isolator with the unionfs strategy")

	if err := strategy.changeRootHandler.close(); err != nil {
		return errors.Wrap(err, "Error closing change root")
	}

	if err := strategy.procMounter.unmount(); err != nil {
		return errors.Wrap(err, "Error unmounting proc")
	}

	if err := strategy.devMounter.unmount(); err != nil {
		return errors.Wrap(err, "Error unmounting dev")
	}

	return strategy.changeIsolator.cleanupInsideNS()
}

func (strategy *changeIsolatorUnionFsStrategy) cleanupOutsideNS() error {
	log.Tracef("[Outside NS] Cleaning up change isolator with the unionfs strategy")

	if err := strategy.unionFsHandler.unmount(); err != nil {
		return errors.Wrap(err, "Error unmounting unionfs")
	}

	if err := strategy.unionFsHandler.cleanup(); err != nil {
		return errors.Wrap(err, "Error cleaning up unionfs")
	}

	return strategy.changeIsolator.cleanupOutsideNS()
}
