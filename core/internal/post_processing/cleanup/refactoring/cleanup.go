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

	stageToPipelineMergeOptions := dirmerger.NewMergeOptions()
	stageToPipelineMergeOptions.BlockOverwrite = false
	stageToPipelineMergeOptions.MergeMetaFilesFolder = false

	for _, stage := range pipeline.Stages {
		if err := CleanupStage(stage); err != nil {
			return errorx.InternalError.Wrap(err, "Error cleaning up stage")
		}

		// merge stage to target dir and allow overwrites
		if err := dirmerger.MergeFolders(
			[]string{stage.ChangeCaptureLocation},
			targetFolder, stageToPipelineMergeOptions,
		); err != nil {
			return errorx.InternalError.Wrap(err, "failed to cleanup stage")
		}
	}

	pipelineToOutputMergeOptions := dirmerger.NewMergeOptions()
	pipelineToOutputMergeOptions.BlockOverwrite = false
	pipelineToOutputMergeOptions.DeleteEmptyFolders = true
	pipelineToOutputMergeOptions.DeleteMarkedAsDeletedPaths = false

	if err := dirmerger.MergeFolders(
		[]string{pipeline.ChangeCaptureLocation},
		targetFolder, pipelineToOutputMergeOptions,
	); err != nil {
		return errorx.InternalError.Wrap(err, "failed to cleanup pipeline")
	}

	if err := os.RemoveAll(filepath.Join(pipeline.ChangeCaptureLocation, "tmp")); err != nil {
		return errorx.InternalError.Wrap(err, "failed to remove temporary changes directory")
	}

	return nil
}

func CleanupStage(stage *refactoringpipelinemodel.Stage) error {
	stepToStageMergeOptions := dirmerger.NewMergeOptions()
	stepToStageMergeOptions.BlockOverwrite = true
	stepToStageMergeOptions.MergeMetaFilesFolder = false

	sourceFolders := make([]string, 0)
	for _, step := range stage.Steps {
		sourceFolders = append(sourceFolders, step.ChangeCaptureLocation)
	}

	// merge steps in stage and prevent overwrites
	if err := dirmerger.MergeFolders(sourceFolders, stage.ChangeCaptureLocation, stepToStageMergeOptions); err != nil {
		return errorx.InternalError.Wrap(err, "failed to cleanup stage")
	}

	return nil
}
