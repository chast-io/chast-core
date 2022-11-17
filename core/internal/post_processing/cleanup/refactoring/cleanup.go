package refactoringpipelinecleanup

import (
	"os"
	"path/filepath"

	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	"chast.io/core/internal/post_processing/merger/pkg/dirmerger"
	"github.com/joomcode/errorx"
)

func CleanupPipeline(pipeline *refactoringpipelinemodel.Pipeline) error {
	targetFolder := pipeline.ChangeCaptureLocation

	for _, stage := range pipeline.Stages {
		if err := CleanupStage(stage); err != nil {
			return errorx.InternalError.Wrap(err, "Error cleaning up stage")
		}

		if err := cleanupDeletedPaths(stage.ChangeCaptureLocation); err != nil {
			return errorx.InternalError.Wrap(err, "Error cleaning up deleted paths in stage")
		}

		// merge stage to target dir and allow overwrites
		if err := dirmerger.MergeFolders([]string{stage.ChangeCaptureLocation}, targetFolder, false); err != nil {
			return errorx.InternalError.Wrap(err, "failed to cleanup stage")
		}
	}

	if err := cleanupDeletedPaths(targetFolder); err != nil {
		return errorx.InternalError.Wrap(err, "Error cleaning up deleted paths in pipeline")
	}

	if err := dirmerger.MergeDeletions(targetFolder); err != nil {
		return errorx.InternalError.Wrap(err, "Error merging deletions")
	}

	if err := os.RemoveAll(filepath.Join(pipeline.ChangeCaptureLocation, "tmp")); err != nil {
		return errorx.InternalError.Wrap(err, "failed to remove temporary changes directory")
	}

	return nil
}

func CleanupStage(stage *refactoringpipelinemodel.Stage) error {
	sourceFolders := make([]string, 0)
	for _, step := range stage.Steps {
		sourceFolders = append(sourceFolders, step.ChangeCaptureLocation)
	}
	// merge steps in stage and prevent overwrites
	if err := dirmerger.MergeFolders(sourceFolders, stage.ChangeCaptureLocation, true); err != nil {
		return errorx.InternalError.Wrap(err, "failed to cleanup stage")
	}

	return nil
}

func cleanupDeletedPaths(sourcePath string) error {
	removedPathsMetaFolder := filepath.Join(sourcePath, ".unionfs-fuse")

	// merge removed folder with step folder. Overwrite-block should not matter as everything is in the same union.
	if err := dirmerger.MergeFolders([]string{removedPathsMetaFolder}, sourcePath, true); err != nil {
		return errorx.InternalError.Wrap(err, "failed to cleanup step")
	}

	return nil
}
