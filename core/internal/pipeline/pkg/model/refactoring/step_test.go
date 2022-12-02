package refactoringpipelinemodel_test

import (
	"path/filepath"
	"strconv"
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

func stepDummyStep(nr int) *uut.Step {
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
		step := stepDummyStep(1)

		executionGroup.AddStep(step)

		expectedLocation := filepath.Join(pipeline.OperationLocation, step.UUID)
		if step.OperationLocation != expectedLocation {
			t.Errorf("Expected step operation location to be '%s', but was '%s'", expectedLocation, step.OperationLocation)
		}
	})

	t.Run("should set change capture location", func(t *testing.T) {
		t.Parallel()

		executionGroup, pipeline := stepDummyExecutionGroupWithPipeline()
		step := stepDummyStep(1)

		executionGroup.AddStep(step)

		expectedLocation := filepath.Join(pipeline.GetTempChangeCaptureLocation(), step.UUID)
		if step.ChangeCaptureLocation != expectedLocation {
			t.Errorf("Expected step change capture location to be '%s', but was '%s'", expectedLocation, step.ChangeCaptureLocation)
		}
	})

	t.Run("should set pipeline", func(t *testing.T) {
		t.Parallel()

		executionGroup, pipeline := stepDummyExecutionGroupWithPipeline()
		step := stepDummyStep(1)

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

	step := stepDummyStep(1)

	t.Run("should add dependency", func(t *testing.T) {
		t.Parallel()
		dependency := stepDummyStep(2)
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
		dependency := stepDummyStep(2)
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

		step := stepDummyStep(1)
		if !step.IsFinalStep() {
			t.Errorf("Expected step to be final step, but was not")
		}
	})

	t.Run("should return false if step has dependents", func(t *testing.T) {
		t.Parallel()

		step := stepDummyStep(1)
		dependent := stepDummyStep(2)
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

		step := stepDummyStep(1)
		step.ChangeCaptureLocation = "changeCaptureLocation"

		if step.GetFinalChangesLocation() != step.ChangeCaptureLocation+"-final" {
			t.Errorf("Expected final changes location to be '%s', but was '%s'", step.ChangeCaptureLocation+"-final", step.GetFinalChangesLocation())
		}
	})
}

// endregion

// region GetMergedPreviousChangesLocation

func TestStep_GetPreviousChangesLocation(t *testing.T) {
	t.Parallel()

	t.Run("should return change previous changes location", func(t *testing.T) {
		t.Parallel()

		step := stepDummyStep(1)
		step.ChangeCaptureLocation = "changeCaptureLocation"

		if step.GetMergedPreviousChangesLocation() != step.ChangeCaptureLocation+"-prev" {
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

		step := stepDummyStep(1)
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

// region GetPreviousChangeCaptureLocations

func TestStep_GetPreviousChangeCaptureLocations(t *testing.T) {
	t.Parallel()

	t.Run("should return empty slice if no previous steps", func(t *testing.T) {
		t.Parallel()

		step := stepDummyStep(1)
		if len(step.GetPreviousChangeCaptureLocations()) != 0 {
			t.Errorf("Expected previous change capture locations to be empty, but was not")
		}
	})

	t.Run("should return previous change capture locations", func(t *testing.T) {
		t.Parallel()

		// step1 <- step2 <- step3 <- stepFinal
		// step4 <- step5 <- step6 <- stepFinal
		// step7 <- step8 <- step9

		step1 := stepDummyStep(1)
		step2 := stepDummyStep(2)
		step3 := stepDummyStep(3)
		step4 := stepDummyStep(4)
		step5 := stepDummyStep(5)
		step6 := stepDummyStep(6)
		step7 := stepDummyStep(7)
		step8 := stepDummyStep(8)
		step9 := stepDummyStep(9)
		stepFinal := stepDummyStep(10)

		step2.AddDependency(step1)
		step3.AddDependency(step2)

		step5.AddDependency(step4)
		step6.AddDependency(step5)

		step8.AddDependency(step7)
		step9.AddDependency(step8)

		stepFinal.AddDependency(step3)
		stepFinal.AddDependency(step6)

		group1 := stepDummyExecutionGroup()
		group2 := stepDummyExecutionGroup()
		group3 := stepDummyExecutionGroup()
		group4 := stepDummyExecutionGroup()

		group1.AddStep(step1)
		group1.AddStep(step4)
		group1.AddStep(step7)

		group2.AddStep(step2)
		group2.AddStep(step5)
		group2.AddStep(step8)

		group3.AddStep(step3)
		group3.AddStep(step6)
		group3.AddStep(step9)

		group4.AddStep(stepFinal)

		pipeline := stepDummyPipeline()
		pipeline.AddExecutionGroup(group1)
		pipeline.AddExecutionGroup(group2)
		pipeline.AddExecutionGroup(group3)
		pipeline.AddExecutionGroup(group4)

		if len(stepFinal.GetPreviousChangeCaptureLocations()) != 6 {
			t.Fatalf("Expected previous change capture locations to be 6, but was %d", len(stepFinal.GetPreviousChangeCaptureLocations()))
		}

		if stepFinal.GetPreviousChangeCaptureLocations()[5] != step3.ChangeCaptureLocation {
			t.Errorf("Expected previous change capture location to be '%s', but was '%s'", step3.ChangeCaptureLocation, stepFinal.GetPreviousChangeCaptureLocations()[5])
		}

		if stepFinal.GetPreviousChangeCaptureLocations()[4] != step6.ChangeCaptureLocation {
			t.Errorf("Expected previous change capture location to be '%s', but was '%s'", step6.ChangeCaptureLocation, stepFinal.GetPreviousChangeCaptureLocations()[4])
		}

		if stepFinal.GetPreviousChangeCaptureLocations()[3] != step2.ChangeCaptureLocation {
			t.Errorf("Expected previous change capture location to be '%s', but was '%s'", step2.ChangeCaptureLocation, stepFinal.GetPreviousChangeCaptureLocations()[3])
		}

		if stepFinal.GetPreviousChangeCaptureLocations()[2] != step5.ChangeCaptureLocation {
			t.Errorf("Expected previous change capture location to be '%s', but was '%s'", step5.ChangeCaptureLocation, stepFinal.GetPreviousChangeCaptureLocations()[2])
		}

		if stepFinal.GetPreviousChangeCaptureLocations()[1] != step1.ChangeCaptureLocation {
			t.Errorf("Expected previous change capture location to be '%s', but was '%s'", step1.ChangeCaptureLocation, stepFinal.GetPreviousChangeCaptureLocations()[1])
		}

		if stepFinal.GetPreviousChangeCaptureLocations()[0] != step4.ChangeCaptureLocation {
			t.Errorf("Expected previous change capture location to be '%s', but was '%s'", step4.ChangeCaptureLocation, stepFinal.GetPreviousChangeCaptureLocations()[0])
		}
	})
}

// endregion
