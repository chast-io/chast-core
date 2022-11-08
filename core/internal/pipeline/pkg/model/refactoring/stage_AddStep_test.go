package refactoringpipelinemodel

import (
	"path/filepath"
	"testing"
)

func TestStage_AddStep(t *testing.T) {
	testNoPipelineSet(t)
	testPipelineSet(t)
}

func testNoPipelineSet(t *testing.T) bool {
	return t.Run("no pipeline set", func(t *testing.T) {
		t.Parallel()
		t.Run("should not add nil step", func(t *testing.T) {
			t.Parallel()
			actualStage := dummyStage()
			actualStage.AddStep(nil)
			if len(actualStage.Steps) != 0 {
				t.Errorf("Expected stage to have 0 steps, but had %d", len(actualStage.Steps))
			}
		})

		t.Run("should not set stage locations", func(t *testing.T) {
			t.Parallel()
			actualStage := dummyStage()
			step := dummyStep()

			actualStage.AddStep(step)

			t.Run("should set change capture location", func(t *testing.T) {
				t.Parallel()
				if actualStage.ChangeCaptureLocation != "" {
					t.Errorf("Expected stage ChangeCaptureLocation to be empty, but was '%s'", actualStage.ChangeCaptureLocation)
				}
			})

			t.Run("should set operation location", func(t *testing.T) {
				t.Parallel()
				if actualStage.OperationLocation != "" {
					t.Errorf("Expected stage OperationLocation to be empty, but was '%s'", actualStage.OperationLocation)
				}
			})

		})

		t.Run("should not set step locations", func(t *testing.T) {
			t.Parallel()
			actualStage := dummyStage()
			step := dummyStep()

			actualStage.AddStep(step)

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
	})
}

func testPipelineSet(t *testing.T) bool {
	return t.Run("pipeline set", func(t *testing.T) {
		t.Parallel()
		t.Run("should not add nil step", func(t *testing.T) {
			t.Parallel()
			actualStage, _ := dummyStageWithPipeline()
			actualStage.AddStep(nil)
			if len(actualStage.Steps) != 0 {
				t.Errorf("Expected stage to have 0 steps, but had %d", len(actualStage.Steps))
			}
		})

		t.Run("should set stage locations", func(t *testing.T) {
			t.Parallel()
			actualStage, pipeline := dummyStageWithPipeline()
			step := dummyStep()

			actualStage.AddStep(step)

			t.Run("should set change capture location", func(t *testing.T) {
				t.Parallel()
				expectedChangeCaptureLocation := filepath.Join(pipeline.ChangeCaptureLocation, "tmp", actualStage.UUID)
				if actualStage.ChangeCaptureLocation != expectedChangeCaptureLocation {
					t.Errorf("Expected stage ChangeCaptureLocation to be '%s', but was '%s'", expectedChangeCaptureLocation, actualStage.ChangeCaptureLocation)
				}
			})

			t.Run("should set operation location", func(t *testing.T) {
				t.Parallel()
				expectedOperationLocation := filepath.Join(pipeline.OperationLocation, actualStage.UUID)
				if actualStage.OperationLocation != expectedOperationLocation {
					t.Errorf("Expected stage OperationLocation to be '%s', but was '%s'", expectedOperationLocation, actualStage.OperationLocation)
				}
			})

		})

		t.Run("should set step locations", func(t *testing.T) {
			t.Parallel()
			actualStage, pipeline := dummyStageWithPipeline()
			step := dummyStep()

			actualStage.AddStep(step)

			t.Run("should set change capture location", func(t *testing.T) {
				t.Parallel()
				expectedChangeCaptureLocation := filepath.Join(pipeline.ChangeCaptureLocation, "tmp", actualStage.UUID, step.UUID)
				if step.ChangeCaptureLocation != expectedChangeCaptureLocation {
					t.Errorf("Expected step ChangeCaptureLocation to be '%s', but was '%s'", expectedChangeCaptureLocation, step.ChangeCaptureLocation)
				}
			})

			t.Run("should set operation location", func(t *testing.T) {
				t.Parallel()
				expectedChangeCaptureLocation := filepath.Join(pipeline.OperationLocation, actualStage.UUID, step.UUID)
				if step.OperationLocation != expectedChangeCaptureLocation {
					t.Errorf("Expected step OperationLocation to be '%s', but was '%s'", expectedChangeCaptureLocation, step.OperationLocation)
				}
			})
		})
	})
}
