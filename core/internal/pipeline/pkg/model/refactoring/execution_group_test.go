package refactoringpipelinemodel_test

import (
	"path/filepath"
	"testing"

	uut "chast.io/core/internal/pipeline/pkg/model/refactoring"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

// region Helpers
func executionGroupDummyPipeline() *uut.Pipeline {
	return uut.NewPipeline("/operationLocation", "/changeCaptureLocation", "/rootFileSystemLocation")
}

func executionGroupDummyExecutionGroup() *uut.ExecutionGroup {
	return uut.NewExecutionGroup()
}

func executionGroupDummyExecutionGroupWithPipeline() (*uut.ExecutionGroup, *uut.Pipeline) {
	executionGroup := executionGroupDummyExecutionGroup()
	pipeline := executionGroupDummyPipeline()
	pipeline.AddExecutionGroup(executionGroup)

	return executionGroup, pipeline
}

func executionGroupDummyStep() *uut.Step {
	runModel := &refactoring.SingleRunModel{
		Run: &refactoring.Run{
			ID:                 "runId",
			Dependencies:       make([]*refactoring.Run, 0),
			SupportedLanguages: []string{"java"},
			Docker:             &refactoring.Docker{},          //nolint:exhaustruct // not required for test
			Local:              &refactoring.Local{},           //nolint:exhaustruct // not required for test
			Command:            &refactoring.Command{},         //nolint:exhaustruct // not required for test
			ChangeLocations:    &refactoring.ChangeLocations{}, //nolint:exhaustruct // not required for test
		},
	}

	return uut.NewStep(runModel)
}

// endregion

// region AddStep
func TestExecutionGroup_NewExecutionGroup(t *testing.T) {
	t.Parallel()

	actualExecutionGroup := uut.NewExecutionGroup()

	t.Run("should set initial steps size", func(t *testing.T) {
		t.Parallel()
		if len(actualExecutionGroup.Steps) != 0 {
			t.Errorf("Expected executionGroup to have 0 steps, but had %d", len(actualExecutionGroup.Steps))
		}
	})
}

// endregion

// region AddStep
func TestExecutionGroup_AddStep(t *testing.T) {
	t.Parallel()

	testNoPipelineSet(t)
	testPipelineSet(t)
}

func testNoPipelineSet(t *testing.T) bool {
	t.Helper()

	return t.Run("no Pipeline set", func(t *testing.T) {
		t.Parallel()
		t.Run("should not add nil step", func(t *testing.T) {
			t.Parallel()
			actualExecutionGroup := executionGroupDummyExecutionGroup()
			actualExecutionGroup.AddStep(nil)
			if len(actualExecutionGroup.Steps) != 0 {
				t.Errorf("Expected executionGroup to have 0 steps, but had %d", len(actualExecutionGroup.Steps))
			}
		})

		t.Run("should not set step locations", func(t *testing.T) {
			t.Parallel()
			actualExecutionGroup := executionGroupDummyExecutionGroup()
			step := executionGroupDummyStep()

			actualExecutionGroup.AddStep(step)

			t.Run("should not set step change capture location", func(t *testing.T) {
				t.Parallel()
				if step.ChangeCaptureLocation != "" {
					t.Errorf("Expected step ChangeCaptureLocation to be empty, but was '%s'", step.ChangeCaptureLocation)
				}
			})

			t.Run("should not set step operation location", func(t *testing.T) {
				t.Parallel()
				if step.OperationLocation != "" {
					t.Errorf("Expected step OperationLocation to be empty, but was '%s'", step.OperationLocation)
				}
			})
		})

		t.Run("should set step locations after setting pipeline", func(t *testing.T) {
			t.Parallel()
			actualExecutionGroup := executionGroupDummyExecutionGroup()
			step := executionGroupDummyStep()

			actualExecutionGroup.AddStep(step)

			pipeline := executionGroupDummyPipeline()
			pipeline.AddExecutionGroup(actualExecutionGroup)

			t.Run("should set change capture location", func(t *testing.T) {
				t.Parallel()
				expectedChangeCaptureLocation := filepath.Join(pipeline.GetTempChangeCaptureLocation(), step.UUID)
				if step.ChangeCaptureLocation != expectedChangeCaptureLocation {
					t.Errorf("Expected step ChangeCaptureLocation to be '%s', but was '%s'", expectedChangeCaptureLocation, step.ChangeCaptureLocation)
				}
			})

			t.Run("should set operation location", func(t *testing.T) {
				t.Parallel()
				expectedChangeCaptureLocation := filepath.Join(pipeline.OperationLocation, step.UUID)
				if step.OperationLocation != expectedChangeCaptureLocation {
					t.Errorf("Expected step OperationLocation to be '%s', but was '%s'", expectedChangeCaptureLocation, step.OperationLocation)
				}
			})
		})
	})
}

func testPipelineSet(t *testing.T) bool {
	t.Helper()

	return t.Run("Pipeline set", func(t *testing.T) {
		t.Parallel()
		t.Run("should not add nil step", func(t *testing.T) {
			t.Parallel()
			actualExecutionGroup, _ := executionGroupDummyExecutionGroupWithPipeline()
			actualExecutionGroup.AddStep(nil)
			if len(actualExecutionGroup.Steps) != 0 {
				t.Errorf("Expected executionGroup to have 0 steps, but had %d", len(actualExecutionGroup.Steps))
			}
		})

		t.Run("should set step locations", func(t *testing.T) {
			t.Parallel()
			actualExecutionGroup, pipeline := executionGroupDummyExecutionGroupWithPipeline()
			step := executionGroupDummyStep()

			actualExecutionGroup.AddStep(step)

			t.Run("should set change capture location", func(t *testing.T) {
				t.Parallel()
				expectedChangeCaptureLocation := filepath.Join(pipeline.GetTempChangeCaptureLocation(), step.UUID)
				if step.ChangeCaptureLocation != expectedChangeCaptureLocation {
					t.Errorf("Expected step ChangeCaptureLocation to be '%s', but was '%s'", expectedChangeCaptureLocation, step.ChangeCaptureLocation)
				}
			})

			t.Run("should set operation location", func(t *testing.T) {
				t.Parallel()
				expectedChangeCaptureLocation := filepath.Join(pipeline.OperationLocation, step.UUID)
				if step.OperationLocation != expectedChangeCaptureLocation {
					t.Errorf("Expected step OperationLocation to be '%s', but was '%s'", expectedChangeCaptureLocation, step.OperationLocation)
				}
			})
		})
	})
}

// endregion
