package pipelinepostprocessor

import (
	"chast.io/core/internal/internal_util/collection"
	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	refactoringpipelinecleanup "chast.io/core/internal/post_processing/cleanup/pkg/refactoring"
	"chast.io/core/internal/post_processing/merger/pkg/dirmerger"
	"chast.io/core/internal/post_processing/merger/pkg/mergeoptions"
	"github.com/joomcode/errorx"
)

func Process(pipeline *refactoringpipelinemodel.Pipeline) error {
	if pipeline == nil {
		return errorx.IllegalArgument.New("pipeline must not be nil")
	}

	// Merge all changes from the final steps into the change capture location of the pipeline.
	// No overwrites must happen
	options := mergeoptions.NewMergeOptions()
	options.BlockOverwrite = true
	options.CopyMode = false
	options.MergeMetaFilesFolder = true
	options.DeleteMarkedAsDeletedPaths = false
	options.DeleteEmptyFolders = false

	targetFolder := pipeline.GetFinalChangeCaptureLocation()

	cumulatedErrors := make([]error, 0)

	mergeEntities := collection.Map(pipeline.GetFinalSteps(),
		func(step *refactoringpipelinemodel.Step) dirmerger.MergeEntity {
			return dirmerger.NewMergeEntity(step.GetFinalChangesLocation(), step.ChangeFilteringLocations())
		})

	if err := dirmerger.MergeFolders(
		mergeEntities,
		targetFolder,
		options,
	); err != nil {
		cumulatedErrors = append(cumulatedErrors,
			errorx.InternalError.Wrap(err, "failed to cleanup stage"))
	}

	if len(cumulatedErrors) > 0 {
		return errorx.WrapMany( //nolint:wrapcheck // errorx.WrapMany is a wrapper
			errorx.InternalError, "failed to cleanup pipeline", cumulatedErrors...)
	}

	if err := refactoringpipelinecleanup.CleanupPipeline(pipeline); err != nil {
		return errorx.InternalError.Wrap(err, "Failed to cleanup pipeline")
	}

	return nil
}
