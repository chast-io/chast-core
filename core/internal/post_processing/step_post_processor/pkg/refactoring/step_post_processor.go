package steppostprocessor

import (
	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
	refactoringpipelinecleanup "chast.io/core/internal/post_processing/cleanup/pkg/refactoring"
	"chast.io/core/internal/post_processing/merger/pkg/dirmerger"
	"chast.io/core/internal/post_processing/merger/pkg/mergeoptions"
	"github.com/joomcode/errorx"
)

func Process(step *refactoringpipelinemodel.Step) error {
	if step == nil {
		return errorx.IllegalArgument.New("step must not be nil")
	}

	if err := createCopyOfChangesAndFlattenMetaFiles(step); err != nil {
		return err
	}

	if err := filterAndMovePreviousChangesToFinalLocation(step); err != nil {
		return err
	}

	if err := mergeChangedWithPreviousChanges(step); err != nil {
		return err
	}

	if err := publishChangesToDependents(step); err != nil {
		return err
	}

	if err := refactoringpipelinecleanup.CleanupStep(step); err != nil {
		return errorx.InternalError.Wrap(err, "Error running cleanup")
	}

	return nil
}

func createCopyOfChangesAndFlattenMetaFiles(step *refactoringpipelinemodel.Step) error {
	// Files are copied to a staging area and meta files are flattened
	// The original files remain untouched for the unionfs to work correctly.
	// During this operation, overwrites can take place due to the order of moving/copying the files.
	// Empty folders are kept until the end, because they may be intended
	options := mergeoptions.NewMergeOptions()
	options.BlockOverwrite = false
	options.CopyMode = true
	options.MergeMetaFilesFolder = true
	options.DeleteMarkedAsDeletedPaths = false
	options.DeleteEmptyFolders = false

	mergeEntities := []dirmerger.MergeEntity{
		dirmerger.NewMergeEntity(step.ChangeCaptureLocation, step.ChangeFilteringLocations()),
	}

	if err := dirmerger.MergeFolders(
		mergeEntities,
		step.GetChangesStagingLocation(),
		options,
	); err != nil {
		return errorx.InternalError.Wrap(err, "failed to create copy of changes and flatten meta files")
	}

	return nil
}

func filterAndMovePreviousChangesToFinalLocation(step *refactoringpipelinemodel.Step) error {
	// The changes are filtered and moved to the final location
	// The files should be written to a new location so no files should be overwritten.
	options := mergeoptions.NewMergeOptions()
	options.BlockOverwrite = true
	options.CopyMode = false
	options.DeleteMarkedAsDeletedPaths = false
	options.DeleteEmptyFolders = false

	mergeEntities := []dirmerger.MergeEntity{
		dirmerger.NewMergeEntity(step.GetMergedPreviousChangesLocation(), step.ChangeFilteringLocations()),
	}

	if err := dirmerger.MergeFolders(
		mergeEntities,
		step.GetFinalChangesLocation(),
		options,
	); err != nil {
		return errorx.InternalError.Wrap(err, "failed to filter and move previous changes to final location")
	}

	return nil
}

func mergeChangedWithPreviousChanges(step *refactoringpipelinemodel.Step) error {
	// The changes are filtered and moved to the final location
	// Here files can be overwritten, as the later step has precedence over the previous step.
	options := mergeoptions.NewMergeOptions()
	options.BlockOverwrite = false
	options.CopyMode = false
	options.DeleteMarkedAsDeletedPaths = false
	options.DeleteEmptyFolders = false

	mergeEntities := []dirmerger.MergeEntity{
		dirmerger.NewMergeEntity(step.GetChangesStagingLocation(), step.ChangeFilteringLocations()),
	}

	if err := dirmerger.MergeFolders(
		mergeEntities,
		step.GetFinalChangesLocation(),
		options,
	); err != nil {
		return errorx.InternalError.Wrap(err, "failed to filter and move previous changes to final location")
	}

	return nil
}

func publishChangesToDependents(step *refactoringpipelinemodel.Step) error {
	mergeEntities := []dirmerger.MergeEntity{
		dirmerger.NewMergeEntity(step.GetFinalChangesLocation(), nil),
	}

	cumulatedErrors := make([]error, 0)

	// The changes are moved to the dependent location
	// (if there exist multiple, files are copied for all of them but the last one)
	// Files should never overwrite existing files, which can happen if the dependent has multiple dependencies.
	baseOptions := mergeoptions.NewMergeOptions()
	baseOptions.BlockOverwrite = true
	baseOptions.DeleteEmptyFolders = false
	baseOptions.DeleteMarkedAsDeletedPaths = false

	for index, dependent := range step.Dependents {
		options := *baseOptions

		if index == len(step.Dependents)-1 {
			// last dependent, move changes instead of copying
			options.CopyMode = false
		} else {
			options.CopyMode = true
		}

		if err := dirmerger.MergeFolders(
			mergeEntities,
			dependent.GetMergedPreviousChangesLocation(),
			&options,
		); err != nil {
			cumulatedErrors = append(cumulatedErrors,
				errorx.InternalError.Wrap(err, "failed to cleanup stage"))
		}
	}

	if len(cumulatedErrors) > 0 {
		return errorx.WrapMany( //nolint:wrapcheck // errorx.WrapMany is a wrapper
			errorx.InternalError, "failed to cleanup step", cumulatedErrors...)
	}

	return nil
}
