package refactoringpipelinecleanup

import (
	"os"
	"path/filepath"

	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	"chast.io/core/pkg/util/fs/folder"
	"github.com/pkg/errors"
)

func Cleanup(pipeline *refactoringpipelinemodel.Pipeline) error {
	targetFolder := pipeline.ChangeCaptureFolder

	for _, stage := range pipeline.Stages {
		sourceFolders := make([]string, 0)
		for _, step := range stage.Steps {
			sourceFolders = append(sourceFolders, step.ChangeCaptureFolder)
		}
		// merge steps in stage and prevent overwrites
		stageTargetFolder := filepath.Join(targetFolder, stage.UUID)
		if err := folder.MergeFolders(sourceFolders, stageTargetFolder, true); err != nil {
			return errors.Wrap(err, "failed to cleanup pipeline")
		}

		// merge stage to target dir and allow overwrites
		if err := folder.MergeFolders([]string{stageTargetFolder}, targetFolder, false); err != nil {
			return errors.Wrap(err, "failed to cleanup pipeline")
		}
	}

	if err := os.RemoveAll(filepath.Join(pipeline.ChangeCaptureFolder, "tmp")); err != nil {
		return errors.Wrap(err, "failed to remove temporary changes directory")
	}

	return nil
}
