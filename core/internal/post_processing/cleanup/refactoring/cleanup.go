package refactoringpipelinecleanup

import (
	"chast.io/core/internal/internal_util/collection"
	"os"

	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	"chast.io/core/internal/post_processing/merger/pkg/dirmerger"
	"github.com/joomcode/errorx"
)

func CleanupPipeline(pipeline *refactoringpipelinemodel.Pipeline) error {
	if pipeline == nil {
		return errorx.IllegalArgument.New("pipeline must not be nil")
	}

	for _, group := range pipeline.ExecutionGroups {
		for _, step := range group.Steps {
			if err := cleanupStep(step, true); err != nil {
				return errorx.InternalError.Wrap(err, "failed to cleanup step")
			}
		}
	}

	targetFolder := pipeline.ChangeCaptureLocation

	pipelineFinalizationMergeOptions := dirmerger.NewMergeOptions()
	pipelineFinalizationMergeOptions.BlockOverwrite = true
	pipelineFinalizationMergeOptions.MergeMetaFilesFolder = true
	pipelineFinalizationMergeOptions.DeleteMarkedAsDeletedPaths = false
	pipelineFinalizationMergeOptions.DeleteEmptyFolders = true

	cleanupErrors := make([]error, 0)

	mergeEntities := collection.Map(pipeline.GetFinalSteps(),
		func(step *refactoringpipelinemodel.Step) dirmerger.MergeEntity {
			return dirmerger.NewMergeEntity(step.GetFinalChangesLocation(), step.ChangeFilteringLocations())
		})

	if err := dirmerger.MergeFolders(
		mergeEntities,
		targetFolder,
		pipelineFinalizationMergeOptions,
	); err != nil {
		cleanupErrors = append(cleanupErrors,
			errorx.InternalError.Wrap(err, "failed to cleanup stage"))
	}

	if err := os.RemoveAll(pipeline.GetTempChangeCaptureLocation()); err != nil {
		cleanupErrors = append(cleanupErrors,
			errorx.ExternalError.Wrap(err, "failed to remove temporary pipeline change capture directory"))
	}

	if len(cleanupErrors) > 0 {
		return errorx.WrapMany( //nolint:wrapcheck // errorx.WrapMany is a wrapper
			errorx.InternalError, "failed to cleanup pipeline", cleanupErrors...)
	}

	return nil
}

func CleanupStep(step *refactoringpipelinemodel.Step) error {
	// meta files folder is not merged here, because it is merged in the pipeline cleanup
	// the meta files are used during the pipeline execution
	return cleanupStep(step, false)
}

func cleanupStep(step *refactoringpipelinemodel.Step, mergeMetaFilesFolder bool) error {
	stepFinalizingMergeOptions := dirmerger.NewMergeOptions()
	stepFinalizingMergeOptions.BlockOverwrite = false
	stepFinalizingMergeOptions.MergeMetaFilesFolder = mergeMetaFilesFolder
	stepFinalizingMergeOptions.DeleteEmptyFolders = false

	mergeEntities := []dirmerger.MergeEntity{
		dirmerger.NewMergeEntity(step.GetPreviousChangesLocation(), step.ChangeFilteringLocations()),
		dirmerger.NewMergeEntity(step.ChangeCaptureLocation, step.ChangeFilteringLocations()),
	}

	cleanupErrors := make([]error, 0)

	if err := dirmerger.MergeFolders(
		mergeEntities,
		step.GetFinalChangesLocation(),
		stepFinalizingMergeOptions,
	); err != nil {
		cleanupErrors = append(cleanupErrors,
			errorx.InternalError.Wrap(err, "failed to cleanup stage"))
	}

	if err := os.RemoveAll(step.ChangeCaptureLocation); err != nil {
		cleanupErrors = append(cleanupErrors,
			errorx.ExternalError.Wrap(err, "failed to remove temporary step change capture directory"))
	}

	if err := os.RemoveAll(step.OperationLocation); err != nil {
		cleanupErrors = append(cleanupErrors,
			errorx.ExternalError.Wrap(err, "failed to remove temporary step change capture directory"))
	}

	if err := publishChangesToDependents(step); err != nil {
		cleanupErrors = append(cleanupErrors,
			errorx.ExternalError.Wrap(err, "failed to publish changes to dependents"))
	}

	if len(cleanupErrors) > 0 {
		return errorx.WrapMany( //nolint:wrapcheck // errorx.WrapMany is a wrapper
			errorx.InternalError, "failed to cleanup step", cleanupErrors...)
	}

	return nil
}

func publishChangesToDependents(step *refactoringpipelinemodel.Step) error {
	mergeEntities := []dirmerger.MergeEntity{
		dirmerger.NewMergeEntity(step.GetFinalChangesLocation(), nil),
	}

	cleanupErrors := make([]error, 0)

	for index, dependent := range step.Dependents {
		changePublishingMergeOptions := dirmerger.NewMergeOptions()
		changePublishingMergeOptions.BlockOverwrite = false
		changePublishingMergeOptions.MergeMetaFilesFolder = false
		changePublishingMergeOptions.DeleteEmptyFolders = false

		if index == len(step.Dependents)-1 {
			// last dependent, move changes instead of copying
			changePublishingMergeOptions.CopyMode = false
		} else {
			changePublishingMergeOptions.CopyMode = true
		}

		if err := dirmerger.MergeFolders(
			mergeEntities,
			dependent.GetPreviousChangesLocation(),
			changePublishingMergeOptions,
		); err != nil {
			cleanupErrors = append(cleanupErrors,
				errorx.InternalError.Wrap(err, "failed to cleanup stage"))
		}
	}

	if len(cleanupErrors) > 0 {
		return errorx.WrapMany( //nolint:wrapcheck // errorx.WrapMany is a wrapper
			errorx.InternalError, "failed to cleanup step", cleanupErrors...)
	}

	return nil
}
