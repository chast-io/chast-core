package refactoringpipelinecleanup

import (
	"os"

	refactoringpipelinemodel "chast.io/core/internal/pipeline/pkg/model/refactoring"
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

	if err := os.RemoveAll(pipeline.GetTempChangeCaptureLocation()); err != nil {
		return errorx.ExternalError.Wrap(err, "failed to remove temporary pipeline change capture directory")
	}

	return nil
}

func CleanupStep(step *refactoringpipelinemodel.Step) error {
	return cleanupStep(step, false)
}

func cleanupStep(step *refactoringpipelinemodel.Step, clearChangeCaptureLocations bool) error {
	cumulatedErrors := make([]error, 0)

	if err := os.RemoveAll(step.OperationLocation); err != nil {
		cumulatedErrors = append(cumulatedErrors,
			errorx.ExternalError.Wrap(err, "failed to remove temporary step change capture directory"))
	}

	if err := os.RemoveAll(step.GetMergedPreviousChangesLocation()); err != nil {
		cumulatedErrors = append(cumulatedErrors,
			errorx.ExternalError.Wrap(err, "failed to remove temporary step change capture directory"))
	}

	if err := os.RemoveAll(step.GetChangesStagingLocation()); err != nil {
		cumulatedErrors = append(cumulatedErrors,
			errorx.ExternalError.Wrap(err, "failed to remove temporary step change capture directory"))
	}

	if clearChangeCaptureLocations {
		if err := os.RemoveAll(step.ChangeCaptureLocation); err != nil {
			cumulatedErrors = append(cumulatedErrors,
				errorx.ExternalError.Wrap(err, "failed to remove temporary step change capture directory"))
		}

		if err := os.RemoveAll(step.GetFinalChangesLocation()); err != nil {
			cumulatedErrors = append(cumulatedErrors,
				errorx.ExternalError.Wrap(err, "failed to remove temporary step change capture directory"))
		}
	}

	if len(cumulatedErrors) > 0 {
		return errorx.WrapMany( //nolint:wrapcheck // errorx.WrapMany is a wrapper
			errorx.InternalError, "failed to cleanup step", cumulatedErrors...)
	}

	return nil
}
