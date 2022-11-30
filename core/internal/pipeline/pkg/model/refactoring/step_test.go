package refactoringpipelinemodel_test

import (
	"path/filepath"
	"testing"

	uut "chast.io/core/internal/pipeline/pkg/model/refactoring"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

// region Helpers
func stepDummyPipeline() *uut.Pipeline {
	return uut.NewPipeline("/operationLocation", "/changeCaptureLocation", "/rootFileSystemLocation")
}

func stepDummyExecutionGroup() *uut.ExecutionGroup {
	return uut.NewExecutionGroup()
}

func stepDummyExecutionGroupWithPipeline() (*uut.ExecutionGroup, *uut.Pipeline) {
	executionGroup := executionGroupDummyExecutionGroup()
	pipeline := executionGroupDummyPipeline()
	pipeline.AddExecutionGroup(executionGroup)

	return executionGroup, pipeline
}

func stepDummyStep() *uut.Step {
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

// region NewStep
func TestNewStep(t *testing.T) {
	t.Parallel()

	id := "runID"
	runModel := &refactoring.SingleRunModel{
		Run: &refactoring.Run{
			ID:                 id,
			Dependencies:       make([]*refactoring.Run, 0),
			SupportedLanguages: []string{"java"},
			Docker:             &refactoring.Docker{},          //nolint:exhaustruct // not required for test
			Local:              &refactoring.Local{},           //nolint:exhaustruct // not required for test
			Command:            &refactoring.Command{},         //nolint:exhaustruct // not required for test
			ChangeLocations:    &refactoring.ChangeLocations{}, //nolint:exhaustruct // not required for test
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

// region withPipeline [indirect]

func TestStep_WithPipeline(t *testing.T) {
	t.Parallel()

	t.Run("should set operation location", func(t *testing.T) {
		t.Parallel()

		executionGroup, pipeline := stepDummyExecutionGroupWithPipeline()
		step := stepDummyStep()

		executionGroup.AddStep(step)

		expectedLocation := filepath.Join(pipeline.OperationLocation, step.UUID)
		if step.OperationLocation != expectedLocation {
			t.Errorf("Expected step operation location to be '%s', but was '%s'", expectedLocation, step.OperationLocation)
		}
	})

	t.Run("should set change capture location", func(t *testing.T) {
		t.Parallel()

		executionGroup, pipeline := stepDummyExecutionGroupWithPipeline()
		step := stepDummyStep()

		executionGroup.AddStep(step)

		expectedLocation := filepath.Join(pipeline.GetTempChangeCaptureLocation(), step.UUID)
		if step.ChangeCaptureLocation != expectedLocation {
			t.Errorf("Expected step change capture location to be '%s', but was '%s'", expectedLocation, step.ChangeCaptureLocation)
		}
	})

	t.Run("should set pipeline", func(t *testing.T) {
		t.Parallel()

		executionGroup, pipeline := stepDummyExecutionGroupWithPipeline()
		step := stepDummyStep()

		executionGroup.AddStep(step)

		if step.Pipeline != pipeline {
			t.Errorf("Expected step pipeline to be '%v', but was '%v'", pipeline, step.Pipeline)
		}
	})
}

// endregion

// region AddDependency

func TestStep_AddDependency(t *testing.T) {
	t.Parallel()

	step := stepDummyStep()

	t.Run("should add dependency", func(t *testing.T) {
		t.Parallel()
		dependency := stepDummyStep()
		step.AddDependency(dependency)

		if len(step.Dependencies) != 1 {
			t.Errorf("Expected step to have 1 dependency, but had %d", len(step.Dependencies))
		}

		if step.Dependencies[0] != dependency {
			t.Errorf("Expected step to have dependency '%v', but had '%v'", dependency, step.Dependencies[0])
		}
	})

	t.Run("should add dependent", func(t *testing.T) {
		t.Parallel()
		dependency := stepDummyStep()
		step.AddDependency(dependency)

		if len(dependency.Dependents) != 1 {
			t.Errorf("Expected step to have 1 dependent, but had %d", len(dependency.Dependents))
		}

		if dependency.Dependents[0] != step {
			t.Errorf("Expected step to be dependent of dependency, but was not")
		}
	})
}

// endregion

// region IsFinalStep

func TestStep_IsFinalStep(t *testing.T) {
	t.Parallel()

	t.Run("should return true if step has no dependents", func(t *testing.T) {
		t.Parallel()

		step := stepDummyStep()
		if !step.IsFinalStep() {
			t.Errorf("Expected step to be final step, but was not")
		}
	})

	t.Run("should return false if step has dependents", func(t *testing.T) {
		t.Parallel()

		step := stepDummyStep()
		dependent := stepDummyStep()
		dependent.AddDependency(step)

		if step.IsFinalStep() {
			t.Errorf("Expected step to not be final step, but was")
		}

		if !dependent.IsFinalStep() {
			t.Errorf("Expected dependent to be final step, but was not")
		}
	})
}

// endregion

// region GetFinalChangesLocation

func TestStep_GetFinalChangesLocation(t *testing.T) {
	t.Parallel()

	t.Run("should return change capture location with suffix", func(t *testing.T) {
		t.Parallel()

		step := stepDummyStep()
		step.ChangeCaptureLocation = "changeCaptureLocation"

		if step.GetFinalChangesLocation() != step.ChangeCaptureLocation+"-final" {
			t.Errorf("Expected final changes location to be '%s', but was '%s'", step.ChangeCaptureLocation+"-final", step.GetFinalChangesLocation())
		}
	})
}

// endregion

// region GetPreviousChangesLocation

func TestStep_GetPreviousChangesLocation(t *testing.T) {
	t.Parallel()

	t.Run("should return change previous changes location", func(t *testing.T) {
		t.Parallel()

		step := stepDummyStep()
		step.ChangeCaptureLocation = "changeCaptureLocation"

		if step.GetPreviousChangesLocation() != step.ChangeCaptureLocation+"-prev" {
			t.Errorf("Expected previous changes location to be '%s', but was '%s'", step.ChangeCaptureLocation+"-prev", step.GetFinalChangesLocation())
		}
	})
}

// endregion

// region ChangeFilteringLocations

func TestStep_ChangeFilteringLocations(t *testing.T) {
	t.Parallel()

	t.Run("should return empty slice if no change locations are set", func(t *testing.T) {
		t.Parallel()

		step := stepDummyStep()
		changeLocations := &refactoring.ChangeLocations{
			Exclude: []string{},
			Include: []string{},
		}
		step.RunModel.Run.ChangeLocations = changeLocations

		if step.ChangeFilteringLocations() != changeLocations {
			t.Errorf("Expected change locations to be '%v', but was '%v'", changeLocations, step.ChangeFilteringLocations())
		}
	})
}

// endregion
