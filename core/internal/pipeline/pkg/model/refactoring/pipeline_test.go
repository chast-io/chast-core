package refactoringpipelinemodel_test

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	uut "chast.io/core/internal/pipeline/pkg/model/refactoring"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

// region Helpers
func pipelineDummyPipeline() *uut.Pipeline {
	return uut.NewPipeline("/operationLocation", "/changeCaptureLocation", "/rootFileSystemLocation")
}

func pipelineDummyStage() *uut.Stage {
	return uut.NewStage("test-name")
}

func pipelineDummyStep() *uut.Step {
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

// region NewPipeline
func TestNewPipeline(t *testing.T) {
	t.Parallel()

	operationLocation := "/operationLocation"
	changeCaptureLocation := "/changeCaptureLocation"
	rootFileSystemLocation := "/rootFileSystemLocation"

	actualPipeline := uut.NewPipeline(operationLocation, changeCaptureLocation, rootFileSystemLocation)

	t.Run("should set UUID prefix", func(t *testing.T) {
		t.Parallel()
		if strings.HasPrefix(actualPipeline.UUID, "PIPELINE-") == false {
			t.Errorf("Expected pipeline UUID to start with 'PIPELINE-', but was '%s'", actualPipeline.UUID)
		}
	})

	t.Run("should set correct UUID", func(t *testing.T) {
		t.Parallel()
		if len(actualPipeline.UUID) != len("PIPELINE-")+len("00000000-0000-0000-0000-000000000000") {
			t.Errorf("Expected pipeline UUID to be 36 characters long, but was %d", len(actualPipeline.UUID))
		}
	})

	t.Run("should set operation location", func(t *testing.T) {
		t.Parallel()
		if actualPipeline.OperationLocation != operationLocation {
			t.Errorf("expected operation location to be %s but was %s", "/tmp/chast", actualPipeline.OperationLocation)
		}
	})

	t.Run("should set change capture location", func(t *testing.T) {
		t.Parallel()
		expectedChangeCaptureLocation := filepath.Join(changeCaptureLocation, actualPipeline.UUID)
		if actualPipeline.ChangeCaptureLocation != expectedChangeCaptureLocation {
			t.Errorf("Expected pipeline ChangeCaptureLocation to be '%s', but was '%s'", expectedChangeCaptureLocation, actualPipeline.ChangeCaptureLocation)
		}
	})

	t.Run("should set root file system location", func(t *testing.T) {
		t.Parallel()
		if actualPipeline.RootFileSystemLocation != rootFileSystemLocation {
			t.Errorf("Expected pipeline RootFileSystemLocation to be '%s', but was '%s'", rootFileSystemLocation, actualPipeline.RootFileSystemLocation)
		}
	})

	t.Run("should set initial stages size", func(t *testing.T) {
		t.Parallel()
		if len(actualPipeline.Stages) != 1 {
			t.Errorf("Expected pipeline to have 1 stage, but had %d", len(actualPipeline.Stages))
		}
	})
}

// endregion

// region AddStage

func TestAddStage(t *testing.T) {
	t.Parallel()

	t.Run("should add stage", func(t *testing.T) {
		t.Parallel()
		pipeline := pipelineDummyPipeline()
		stage := pipelineDummyStage()
		pipeline.AddStage(stage)

		if pipeline.Stages[0] != stage {
			t.Errorf("Expected pipeline to have stage %v, but had %v", stage, pipeline.Stages[0])
		}
	})

	t.Run("should set change capture location on stage", func(t *testing.T) {
		t.Parallel()
		pipeline := pipelineDummyPipeline()
		stage := pipelineDummyStage()
		pipeline.AddStage(stage)

		expectedChangeCaptureLocation := filepath.Join(pipeline.ChangeCaptureLocation, "tmp", stage.UUID)
		if stage.ChangeCaptureLocation != expectedChangeCaptureLocation {
			t.Errorf("Expected stage ChangeCaptureLocation to be '%s', but was '%s'", expectedChangeCaptureLocation, stage.ChangeCaptureLocation)
		}
	})

	t.Run("should set operation location on stage", func(t *testing.T) {
		t.Parallel()
		pipeline := pipelineDummyPipeline()
		stage := pipelineDummyStage()
		pipeline.AddStage(stage)

		expectedOperationLocation := filepath.Join(pipeline.OperationLocation, stage.UUID)
		if stage.OperationLocation != expectedOperationLocation {
			t.Errorf("Expected stage OperationLocation to be '%s', but was '%s'", expectedOperationLocation, stage.OperationLocation)
		}
	})

	t.Run("should set prev stage on stage", func(t *testing.T) {
		t.Parallel()
		pipeline := pipelineDummyPipeline()
		stage1 := uut.NewStage("stage1")
		pipeline.AddStage(stage1)
		stage2 := uut.NewStage("stage2")
		pipeline.AddStage(stage2)

		t.Run("should not set prev on first", func(t *testing.T) {
			t.Parallel()
			if len(stage1.GetPrevChangeCaptureLocations()) != 0 {
				t.Errorf("Expected stage1 to have 0 prev change capture locations, but had %d", len(stage1.GetPrevChangeCaptureLocations()))
			}
		})

		t.Run("should set prev on second", func(t *testing.T) {
			t.Parallel()
			if len(stage2.GetPrevChangeCaptureLocations()) != 1 {
				t.Errorf("Expected stage2 to have 1 prev change capture location, but had %d", len(stage2.GetPrevChangeCaptureLocations()))
			}
		})
	})

	t.Run("should set operation location on already added steps", func(t *testing.T) {
		t.Parallel()
		stage := pipelineDummyStage()
		step1 := pipelineDummyStep()
		stage.AddStep(step1)
		step2 := pipelineDummyStep()
		stage.AddStep(step2)

		pipeline := pipelineDummyPipeline()
		pipeline.AddStage(stage)

		expectedOperationLocation1 := filepath.Join(stage.OperationLocation, step1.UUID)
		if step1.OperationLocation != expectedOperationLocation1 {
			t.Errorf("Expected step OperationLocation to be '%s', but was '%s'", expectedOperationLocation1, step1.OperationLocation)
		}

		expectedOperationLocation2 := filepath.Join(stage.OperationLocation, step2.UUID)
		if step2.OperationLocation != expectedOperationLocation2 {
			t.Errorf("Expected step OperationLocation to be '%s', but was '%s'", expectedOperationLocation2, step2.OperationLocation)
		}
	})
}

// endregion

// region GetPrevChangeCaptureLocations

func TestGetPrevChangeCaptureLocations(t *testing.T) {
	t.Parallel()

	t.Run("should return prev change capture locations in correct order", func(t *testing.T) {
		t.Parallel()
		stage1 := pipelineDummyStage()
		stage2 := pipelineDummyStage()
		stage3 := pipelineDummyStage()
		stage4 := pipelineDummyStage()
		pipeline := pipelineDummyPipeline()
		pipeline.AddStage(stage1)
		pipeline.AddStage(stage2)
		pipeline.AddStage(stage3)
		pipeline.AddStage(stage4)

		if len(stage1.GetPrevChangeCaptureLocations()) != 0 {
			t.Errorf("Expected stage1 to have 0 prev change capture locations, but had %d", len(stage1.GetPrevChangeCaptureLocations()))
		}
		if len(stage2.GetPrevChangeCaptureLocations()) != 1 {
			t.Errorf("Expected stage2 to have 1 prev change capture location, but had %d", len(stage2.GetPrevChangeCaptureLocations()))
		}
		if len(stage3.GetPrevChangeCaptureLocations()) != 2 {
			t.Errorf("Expected stage3 to have 2 prev change capture locations, but had %d", len(stage3.GetPrevChangeCaptureLocations()))
		}
		if len(stage4.GetPrevChangeCaptureLocations()) != 3 {
			t.Errorf("Expected stage4 to have 3 prev change capture locations, but had %d", len(stage4.GetPrevChangeCaptureLocations()))
		}

		expected := []string{
			stage1.ChangeCaptureLocation,
			stage2.ChangeCaptureLocation,
			stage3.ChangeCaptureLocation,
		}
		actual := stage4.GetPrevChangeCaptureLocations()

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expected prev change capture locations to be %v, but was %v", expected, actual)
		}
	})
}

// endregion
