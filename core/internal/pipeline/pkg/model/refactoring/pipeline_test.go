package refactoringpipelinemodel_test

import (
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	uut "chast.io/core/internal/pipeline/pkg/model/refactoring"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

// region Helpers
func pipelineDummyPipeline() *uut.Pipeline {
	return uut.NewPipeline("/operationLocation", "/changeCaptureLocation", "/rootFileSystemLocation")
}

func pipelineDummyExecutionGroup() *uut.ExecutionGroup {
	return uut.NewExecutionGroup()
}

func pipelineDummyStep(nr int) *uut.Step {
	runModel := &refactoring.SingleRunModel{
		Run: &refactoring.Run{
			ID:                 "runId" + strconv.Itoa(nr),
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
			t.Errorf("Expected Pipeline UUID to start with 'PIPELINE-', but was '%s'", actualPipeline.UUID)
		}
	})

	t.Run("should set correct UUID", func(t *testing.T) {
		t.Parallel()
		if len(actualPipeline.UUID) != len("PIPELINE-")+len("00000000-0000-0000-0000-000000000000") {
			t.Errorf("Expected Pipeline UUID to be 36 characters long, but was %d", len(actualPipeline.UUID))
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
			t.Errorf("Expected Pipeline ChangeCaptureLocation to be '%s', but was '%s'", expectedChangeCaptureLocation, actualPipeline.ChangeCaptureLocation)
		}
	})

	t.Run("should set root file system location", func(t *testing.T) {
		t.Parallel()
		if actualPipeline.RootFileSystemLocation != rootFileSystemLocation {
			t.Errorf("Expected Pipeline RootFileSystemLocation to be '%s', but was '%s'", rootFileSystemLocation, actualPipeline.RootFileSystemLocation)
		}
	})

	t.Run("should set initial executionGroups size", func(t *testing.T) {
		t.Parallel()
		if len(actualPipeline.ExecutionGroups) != 1 {
			t.Errorf("Expected Pipeline to have 1 executionGroup, but had %d", len(actualPipeline.ExecutionGroups))
		}
	})
}

// endregion

// region AddExecutionGroup

func TestAddExecutionGroup(t *testing.T) {
	t.Parallel()

	t.Run("should add execution group", func(t *testing.T) {
		t.Parallel()

		pipeline := pipelineDummyPipeline()
		executionGroup := pipelineDummyExecutionGroup()
		pipeline.AddExecutionGroup(executionGroup)

		if pipeline.ExecutionGroups[0] != executionGroup {
			t.Errorf("Expected Pipeline to have executionGroup %v, but had %v", executionGroup, pipeline.ExecutionGroups[0])
		}
	})

	t.Run("should set operation location on already added steps", func(t *testing.T) {
		t.Parallel()

		executionGroup := pipelineDummyExecutionGroup()
		step1 := pipelineDummyStep(1)
		executionGroup.AddStep(step1)
		step2 := pipelineDummyStep(2)
		executionGroup.AddStep(step2)

		pipeline := pipelineDummyPipeline()
		pipeline.AddExecutionGroup(executionGroup)

		expectedOperationLocation1 := filepath.Join(pipeline.OperationLocation, step1.UUID)
		if step1.OperationLocation != expectedOperationLocation1 {
			t.Errorf("Expected step OperationLocation to be '%s', but was '%s'", expectedOperationLocation1, step1.OperationLocation)
		}

		expectedOperationLocation2 := filepath.Join(pipeline.OperationLocation, step2.UUID)
		if step2.OperationLocation != expectedOperationLocation2 {
			t.Errorf("Expected step OperationLocation to be '%s', but was '%s'", expectedOperationLocation2, step2.OperationLocation)
		}
	})

	t.Run("should accept multiple groups", func(t *testing.T) {
		t.Parallel()

		pipeline := pipelineDummyPipeline()
		executionGroup1 := pipelineDummyExecutionGroup()
		pipeline.AddExecutionGroup(executionGroup1)
		executionGroup2 := pipelineDummyExecutionGroup()
		pipeline.AddExecutionGroup(executionGroup2)

		if pipeline.ExecutionGroups[0] != executionGroup1 {
			t.Errorf("Expected Pipeline to have executionGroup %v, but had %v", executionGroup1, pipeline.ExecutionGroups[0])
		}
		if pipeline.ExecutionGroups[1] != executionGroup2 {
			t.Errorf("Expected Pipeline to have executionGroup %v, but had %v", executionGroup2, pipeline.ExecutionGroups[1])
		}
	})
}

// endregion

// region GetTempChangeCaptureLocation

func TestGetTempChangeCaptureLocation(t *testing.T) {
	t.Parallel()

	t.Run("should return correct location", func(t *testing.T) {
		t.Parallel()
		pipeline := pipelineDummyPipeline()
		actualLocation := pipeline.GetTempChangeCaptureLocation()
		expectedLocation := filepath.Join(pipeline.ChangeCaptureLocation, "tmp")
		if actualLocation != expectedLocation {
			t.Errorf("Expected location to be '%s', but was '%s'", expectedLocation, actualLocation)
		}
	})
}

// endregion

// region GetFinalSteps

func TestGetFinalSteps(t *testing.T) {
	t.Parallel()

	t.Run("should return correct steps", func(t *testing.T) {
		t.Parallel()

		pipeline := pipelineDummyPipeline()
		executionGroup := pipelineDummyExecutionGroup()
		step1 := pipelineDummyStep(1)
		executionGroup.AddStep(step1)
		step2 := pipelineDummyStep(2)
		executionGroup.AddStep(step2)
		pipeline.AddExecutionGroup(executionGroup)

		actualSteps := pipeline.GetFinalSteps()

		if len(actualSteps) != 2 {
			t.Fatalf("Expected 2 steps, but had %d", len(actualSteps))
		}

		if actualSteps[0] != step1 {
			t.Errorf("Expected step to be %v, but was %v", step1, actualSteps[0])
		}

		if actualSteps[1] != step2 {
			t.Errorf("Expected step to be %v, but was %v", step2, actualSteps[1])
		}
	})

	t.Run("should return correct steps with dependencies", func(t *testing.T) {
		t.Parallel()

		pipeline := pipelineDummyPipeline()
		executionGroup := pipelineDummyExecutionGroup()
		step1 := pipelineDummyStep(1)
		executionGroup.AddStep(step1)
		step2 := pipelineDummyStep(2)
		step2.AddDependency(step1)
		executionGroup.AddStep(step2)
		pipeline.AddExecutionGroup(executionGroup)

		actualSteps := pipeline.GetFinalSteps()

		if len(actualSteps) != 1 {
			t.Fatalf("Expected 1 steps, but had %d", len(actualSteps))
		}

		if actualSteps[0] != step2 {
			t.Errorf("Expected step to be %v, but was %v", step2, actualSteps[0])
		}
	})
}

// endregion
