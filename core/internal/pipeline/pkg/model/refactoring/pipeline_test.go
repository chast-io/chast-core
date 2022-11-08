package refactoringpipelinemodel

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestNewPipeline(t *testing.T) {
	operationLocation := "/operationLocation"
	changeCaptureLocation := "/changeCaptureLocation"
	rootFileSystemLocation := "/rootFileSystemLocation"

	actualPipeline := NewPipeline(operationLocation, changeCaptureLocation, rootFileSystemLocation)

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

// TODO - add tests for AddStage
