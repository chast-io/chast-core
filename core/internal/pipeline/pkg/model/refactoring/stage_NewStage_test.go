package refactoringpipelinemodel

import (
	"strings"
	"testing"
)

func TestStage_NewStage(t *testing.T) {
	name := "test-name"
	actualStage := NewStage(name)

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

// TODO add test for AddStep, GetPrevChangeCaptureLocations
