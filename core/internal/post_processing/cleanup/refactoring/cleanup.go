package refactoringpipelinecleanup

import (
	"os"

	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	"chast.io/core/internal/post_processing/merger/pkg/dirmerger"
	"github.com/joomcode/errorx"
)

func CleanupPipeline(pipeline *refactoringpipelinemodel.Pipeline) error {
	if pipeline == nil {
		return errorx.IllegalArgument.New("pipeline must not be nil")
	}

	targetFolder := pipeline.ChangeCaptureLocation

	stageToPipelineMergeOptions := dirmerger.NewMergeOptions()
	stageToPipelineMergeOptions.BlockOverwrite = false
	stageToPipelineMergeOptions.MergeMetaFilesFolder = true
	stageToPipelineMergeOptions.DeleteMarkedAsDeletedPaths = false
	stageToPipelineMergeOptions.DeleteEmptyFolders = true

	cleanupErrors := make([]error, 0)

	for _, stage := range pipeline.Stages {
		if err := CleanupStage(stage); err != nil {
			cleanupErrors = append(cleanupErrors,
				errorx.InternalError.Wrap(err, "Error cleaning up stage"))
		}

		// merge stage to target dir and allow overwrites
		if err := dirmerger.MergeFolders(
			[]string{stage.ChangeCaptureLocation},
			targetFolder,
			stageToPipelineMergeOptions,
		); err != nil {
			cleanupErrors = append(cleanupErrors,
				errorx.InternalError.Wrap(err, "failed to cleanup stage %s to pipeline %s", stage.UUID, pipeline.UUID))
		}
	}

	if err := os.RemoveAll(pipeline.GetTempChangeCaptureLocation()); err != nil {
		cleanupErrors = append(cleanupErrors,
			errorx.ExternalError.Wrap(err, "failed to remove temporary change capture directory"))
	}

	if err := os.RemoveAll(pipeline.OperationLocation); err != nil {
		cleanupErrors = append(cleanupErrors,
			errorx.ExternalError.Wrap(err, "failed to remove temporary operation directory"))
	}

	if len(cleanupErrors) > 0 {
		return errorx.WrapMany( //nolint:wrapcheck // errorx.WrapMany is a wrapper
			errorx.InternalError, "failed to cleanup pipeline", cleanupErrors...)
	}

	return nil
}

func CleanupStage(stage *refactoringpipelinemodel.Stage) error {
	stepToStageMergeOptions := dirmerger.NewMergeOptions()
	stepToStageMergeOptions.BlockOverwrite = true
	stepToStageMergeOptions.MergeMetaFilesFolder = false
	stepToStageMergeOptions.DeleteEmptyFolders = false

	sourceFolders := make([]string, 0)
	for _, step := range stage.Steps {
		sourceFolders = append(sourceFolders, step.ChangeCaptureLocation)
	}

	cleanupErrors := make([]error, 0)

	// merge steps in stage and prevent overwrites
	if err := dirmerger.MergeFolders(sourceFolders, stage.ChangeCaptureLocation, stepToStageMergeOptions); err != nil {
		cleanupErrors = append(cleanupErrors,
			errorx.InternalError.Wrap(err, "failed to cleanup stage"))
	}

	if err := os.RemoveAll(stage.OperationLocation); err != nil {
		cleanupErrors = append(cleanupErrors,
			errorx.ExternalError.Wrap(err, "failed to remove temporary operation directory"))
	}

	if len(cleanupErrors) > 0 {
		return errorx.WrapMany( //nolint:wrapcheck // errorx.WrapMany is a wrapper
			errorx.InternalError, "failed to cleanup stage", cleanupErrors...)
	}

	return nil
}
