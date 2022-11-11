package refactoringpipelinemodel_test

import (
	"path/filepath"
	"strings"
	"testing"

	uut "chast.io/core/internal/pipeline/pkg/model/refactoring"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

// region Helpers
func stageDummyPipeline() *uut.Pipeline {
	return uut.NewPipeline("/operationLocation", "/changeCaptureLocation", "/rootFileSystemLocation")
}

func stageDummyStage() *uut.Stage {
	return uut.NewStage("test-name")
}

func stageDummyStageWithPipeline() (*uut.Stage, *uut.Pipeline) {
	stage := stageDummyStage()
	pipeline := stageDummyPipeline()
	pipeline.AddStage(stage)

	return stage, pipeline
}

func stageDummyStep() *uut.Step {
	runModel := &refactoring.SingleRunModel{
		Run: &refactoring.Run{
			ID:                 "runId",
			Dependencies:       make([]*refactoring.Run, 0),
			SupportedLanguages: []string{"java"},
			Docker:             &refactoring.Docker{},  //nolint:exhaustruct // not required for test
			Local:              &refactoring.Local{},   //nolint:exhaustruct // not required for test
			Command:            &refactoring.Command{}, //nolint:exhaustruct // not required for test
		},
	}

	return uut.NewStep(runModel)
}

// endregion

// region AddStep
func TestStage_NewStage(t *testing.T) {
	t.Parallel()

	name := "test-name"
	actualStage := uut.NewStage(name)

	t.Run("should set UUID prefix", func(t *testing.T) {
		t.Parallel()
		if strings.HasPrefix(actualStage.UUID, "STAGE-"+name+"-") == false {
			t.Errorf("Expected stage UUID to start with 'STAGE-test-name-', but was '%s'", actualStage.UUID)
		}
	})

	t.Run("should set correct UUID ", func(t *testing.T) {
		t.Parallel()
		if len(actualStage.UUID) != len("STAGE-"+name+"-")+len("00000000-0000-0000-0000-000000000000") {
			t.Errorf("Expected stage UUID to be 36 characters long, but was %d", len(actualStage.UUID))
		}
	})

	t.Run("should not set change capture location", func(t *testing.T) {
		t.Parallel()
		if actualStage.ChangeCaptureLocation != "" {
			t.Errorf("Expected stage ChangeCaptureLocation to be empty, but was '%s'", actualStage.ChangeCaptureLocation)
		}
	})

	t.Run("should not set operation location", func(t *testing.T) {
		t.Parallel()
		if actualStage.OperationLocation != "" {
			t.Errorf("Expected stage OperationLocation to be empty, but was '%s'", actualStage.OperationLocation)
		}
	})

	t.Run("should set initial steps size", func(t *testing.T) {
		t.Parallel()
		if len(actualStage.Steps) != 0 {
			t.Errorf("Expected stage to have 0 steps, but had %d", len(actualStage.Steps))
		}
	})

	t.Run("should has no previous change capture locations", func(t *testing.T) {
		t.Parallel()
		prevChangeCaptureLocations := actualStage.GetPrevChangeCaptureLocations()
		if len(prevChangeCaptureLocations) != 0 {
			t.Errorf("Expected stage to have 0 previous change capture locations, but had %d", len(prevChangeCaptureLocations))
		}
	})
}

// endregion

// region AddStep
func TestStage_AddStep(t *testing.T) {
	t.Parallel()

	testNoPipelineSet(t)
	testPipelineSet(t)
}

func testNoPipelineSet(t *testing.T) bool {
	t.Helper()

	return t.Run("no pipeline set", func(t *testing.T) {
		t.Parallel()
		t.Run("should not add nil step", func(t *testing.T) {
			t.Parallel()
			actualStage := stageDummyStage()
			actualStage.AddStep(nil)
			if len(actualStage.Steps) != 0 {
				t.Errorf("Expected stage to have 0 steps, but had %d", len(actualStage.Steps))
			}
		})

		t.Run("should not set stage locations", func(t *testing.T) {
			t.Parallel()
			actualStage := stageDummyStage()
			step := stageDummyStep()

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
			actualStage := stageDummyStage()
			step := stageDummyStep()

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
	t.Helper()

	return t.Run("pipeline set", func(t *testing.T) {
		t.Parallel()
		t.Run("should not add nil step", func(t *testing.T) {
			t.Parallel()
			actualStage, _ := stageDummyStageWithPipeline()
			actualStage.AddStep(nil)
			if len(actualStage.Steps) != 0 {
				t.Errorf("Expected stage to have 0 steps, but had %d", len(actualStage.Steps))
			}
		})

		t.Run("should set stage locations", func(t *testing.T) {
			t.Parallel()
			actualStage, pipeline := stageDummyStageWithPipeline()
			step := stageDummyStep()

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
			actualStage, pipeline := stageDummyStageWithPipeline()
			step := stageDummyStep()

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

// endregion
