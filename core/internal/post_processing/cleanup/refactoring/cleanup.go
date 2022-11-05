package refactoringpipelinecleanup

import (
	"os"
	"path/filepath"

	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	"chast.io/core/internal/post_processing/merger/pkg/dirmerger"
	"github.com/pkg/errors"
)

func CleanupPipeline(pipeline *refactoringpipelinemodel.Pipeline) error {
	targetFolder := pipeline.ChangeCaptureFolder

	for _, stage := range pipeline.Stages {
		if err := CleanupStage(stage); err != nil {
			return errors.Wrap(err, "Error cleaning up stage")
		}

		if err := cleanupDeletedPaths(stage.ChangeCaptureFolder); err != nil {
			return errors.Wrap(err, "Error cleaning up deleted paths in stage")
		}

		// merge stage to target dir and allow overwrites
		if err := dirmerger.MergeFolders([]string{stage.ChangeCaptureFolder}, targetFolder, false); err != nil {
			return errors.Wrap(err, "failed to cleanup stage")
		}
	}

	if err := cleanupDeletedPaths(targetFolder); err != nil {
		return errors.Wrap(err, "Error cleaning up deleted paths in pipeline")
	}

	if err := dirmerger.MergeDeletions(targetFolder); err != nil {
		return errors.Wrap(err, "Error merging deletions")
	}

	if err := os.RemoveAll(filepath.Join(pipeline.ChangeCaptureFolder, "tmp")); err != nil {
		return errors.Wrap(err, "failed to remove temporary changes directory")
	}

	return nil
}

func CleanupStage(stage *refactoringpipelinemodel.Stage) error {
	sourceFolders := make([]string, 0)
	for _, step := range stage.Steps {
		sourceFolders = append(sourceFolders, step.ChangeCaptureFolder)
	}
	// merge steps in stage and prevent overwrites
	if err := dirmerger.MergeFolders(sourceFolders, stage.ChangeCaptureFolder, true); err != nil {
		return errors.Wrap(err, "failed to cleanup stage")
	}

	return nil
}

func cleanupDeletedPaths(sourcePath string) error {
	removedPathsMetaFolder := filepath.Join(sourcePath, ".unionfs-fuse")

	// merge removed folder with step folder. Overwrite-block should not matter as everything is in the same union.
	if err := dirmerger.MergeFolders([]string{removedPathsMetaFolder}, sourcePath, true); err != nil {
		return errors.Wrap(err, "failed to cleanup step")
	}

	return nil
}
