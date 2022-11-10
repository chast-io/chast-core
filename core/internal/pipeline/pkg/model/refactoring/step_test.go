package refactoringpipelinemodel_test

import (
	"testing"

	uut "chast.io/core/internal/pipeline/pkg/model/refactoring"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

// region NewStep
func TestNewStep(t *testing.T) {
	t.Parallel()

	id := "runID"
	runModel := &refactoring.SingleRunModel{
		Run: &refactoring.Run{
			ID:                 id,
			Dependencies:       make([]*refactoring.Run, 0),
			SupportedLanguages: []string{"java"},
			Docker:             refactoring.Docker{},  //nolint:exhaustruct // not required for test
			Local:              refactoring.Local{},   //nolint:exhaustruct // not required for test
			Command:            refactoring.Command{}, //nolint:exhaustruct // not required for test
		},
	}

	actualStep := uut.NewStep(runModel)

	t.Run("should set UUID from run model", func(t *testing.T) {
		t.Parallel()
		if actualStep.UUID != runModel.Run.GetUUID() {
			t.Errorf("Expected step UUID to be '%s', but was '%s'", runModel.Run.GetUUID(), actualStep.UUID)
		}
	})

	t.Run("should set run model", func(t *testing.T) {
		t.Parallel()
		if actualStep.RunModel != runModel {
			t.Errorf("Expected step run model to be '%v', but was '%v'", runModel, actualStep.RunModel)
		}
	})

	t.Run("should not set change capture location", func(t *testing.T) {
		t.Parallel()
		if actualStep.ChangeCaptureLocation != "" {
			t.Errorf("Expected step change capture location to be empty, but was '%s'", actualStep.ChangeCaptureLocation)
		}
	})

	t.Run("should not set operation location", func(t *testing.T) {
		t.Parallel()
		if actualStep.OperationLocation != "" {
			t.Errorf("Expected step operation location to be empty, but was '%s'", actualStep.OperationLocation)
		}
	})
}

// endregion
