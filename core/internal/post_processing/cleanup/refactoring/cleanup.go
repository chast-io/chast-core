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

	stageMergeOptions := dirmerger.NewMergeOptions()
	stageMergeOptions.BlockOverwrite = false
	stageMergeOptions.MergeMetaFilesFolder = false

	for _, stage := range pipeline.Stages {
		if err := CleanupStage(stage); err != nil {
			return errorx.InternalError.Wrap(err, "Error cleaning up stage")
		}

		// merge stage to target dir and allow overwrites
		if err := dirmerger.MergeFolders([]string{stage.ChangeCaptureLocation}, targetFolder, stageMergeOptions); err != nil {
			return errorx.InternalError.Wrap(err, "failed to cleanup stage")
		}
	}

	pipelineEndMergeOptions := dirmerger.NewMergeOptions()
	pipelineEndMergeOptions.BlockOverwrite = false
	pipelineEndMergeOptions.DeleteEmptyFolders = true
	pipelineEndMergeOptions.DeleteMarkedAsDeletedPaths = false

	if err := dirmerger.MergeFolders(
		[]string{pipeline.ChangeCaptureLocation},
		targetFolder, pipelineEndMergeOptions,
	); err != nil {
		return errorx.InternalError.Wrap(err, "failed to cleanup pipeline")
	}

	if err := os.RemoveAll(filepath.Join(pipeline.ChangeCaptureLocation, "tmp")); err != nil {
		return errorx.InternalError.Wrap(err, "failed to remove temporary changes directory")
	}

	return nil
}

func CleanupStage(stage *refactoringpipelinemodel.Stage) error {
	stepsMergeOptions := dirmerger.NewMergeOptions()
	stepsMergeOptions.BlockOverwrite = true
	stepsMergeOptions.MergeMetaFilesFolder = false

	sourceFolders := make([]string, 0)
	for _, step := range stage.Steps {
		sourceFolders = append(sourceFolders, step.ChangeCaptureLocation)
	}

	// merge steps in stage and prevent overwrites
	if err := dirmerger.MergeFolders(sourceFolders, stage.ChangeCaptureLocation, stepsMergeOptions); err != nil {
		return errorx.InternalError.Wrap(err, "failed to cleanup stage")
	}

	return nil
}
