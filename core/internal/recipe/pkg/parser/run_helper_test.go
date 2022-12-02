package parser_test

import (
	"reflect"
	"testing"

	recipemodel "chast.io/core/internal/recipe/pkg/model"
)

//nolint:gocognit // nested test function
func testRun(t *testing.T,
	run *recipemodel.Run,
	expectedRun recipemodel.Run) {
	t.Helper()

	if run == nil {
		t.Error("Expected run to be set, but was nil")
	}

	t.Run("ID", func(t *testing.T) {
		t.Parallel()

		if run.ID != expectedRun.ID {
			t.Errorf("Expected run ID to be '%v', but was '%v'", expectedRun.ID, run.ID)
		}
	})

	// Dependencies, SupportedExtensions, Flags, (Docker, Local,) Script, ChangeFilteringLocations

	t.Run("Dependencies", func(t *testing.T) {
		t.Parallel()

		if len(run.Dependencies) != len(expectedRun.Dependencies) {
			t.Fatalf("Expected run dependencies to be '%v', but was '%v'", expectedRun.Dependencies, run.Dependencies)
		}

		for i, dependency := range run.Dependencies {
			if dependency != expectedRun.Dependencies[i] {
				t.Errorf("Expected run dependency to be '%v', but was '%v'", expectedRun.Dependencies[i], dependency)
			}
		}
	})

	t.Run("SupportedExtensions", func(t *testing.T) {
		t.Parallel()

		if len(run.SupportedExtensions) != len(expectedRun.SupportedExtensions) {
			t.Fatalf("Expected run supported extensions to be '%v', but was '%v'", expectedRun.SupportedExtensions, run.SupportedExtensions)
		}

		for i, supportedExtension := range run.SupportedExtensions {
			if supportedExtension != expectedRun.SupportedExtensions[i] {
				t.Errorf("Expected run supported extension to be '%v', but was '%v'", expectedRun.SupportedExtensions[i], supportedExtension)
			}
		}
	})

	t.Run("Flags", func(t *testing.T) {
		t.Parallel()

		if len(run.Flags) != len(expectedRun.Flags) {
			t.Fatalf("Expected run flags to be '%v', but was '%v'", expectedRun.Flags, run.Flags)
		}

		for i, flag := range run.Flags {
			if reflect.DeepEqual(flag, expectedRun) {
				t.Errorf("Expected run flag to be '%v', but was '%v'", expectedRun.Flags[i], flag)
			}
		}
	})

	t.Run("Docker", func(t *testing.T) {
		t.Parallel()

		if !reflect.DeepEqual(run.Docker, expectedRun.Docker) {
			t.Errorf("Expected run docker to be '%v', but was '%v'", expectedRun.Docker, run.Docker)
		}
	})

	t.Run("Local", func(t *testing.T) {
		t.Parallel()

		if !reflect.DeepEqual(run.Local, expectedRun.Local) {
			t.Errorf("Expected run local to be '%v', but was '%v'", expectedRun.Local, run.Local)
		}
	})

	t.Run("Script", func(t *testing.T) {
		t.Parallel()

		if !reflect.DeepEqual(run.Script, expectedRun.Script) {
			t.Errorf("Expected run script to be '%v', but was '%v'", expectedRun.Script, run.Script)
		}
	})

	t.Run("IncludeChangeLocations", func(t *testing.T) {
		t.Parallel()

		if len(run.IncludeChangeLocations) != len(expectedRun.IncludeChangeLocations) {
			t.Fatalf("Expected included run change locations to be '%v', but was '%v'", expectedRun.IncludeChangeLocations, run.IncludeChangeLocations)
		}

		for i, changeLocation := range run.IncludeChangeLocations {
			if !reflect.DeepEqual(changeLocation, expectedRun.IncludeChangeLocations[i]) {
				t.Errorf("Expected included run change location to be '%v', but was '%v'", expectedRun.IncludeChangeLocations[i], changeLocation)
			}
		}
	})

	t.Run("ExcludeChangeLocations", func(t *testing.T) {
		t.Parallel()

		if len(run.ExcludeChangeLocations) != len(expectedRun.ExcludeChangeLocations) {
			t.Fatalf("Expected excluded run change locations to be '%v', but was '%v'", expectedRun.ExcludeChangeLocations, run.ExcludeChangeLocations)
		}

		for i, changeLocation := range run.ExcludeChangeLocations {
			if !reflect.DeepEqual(changeLocation, expectedRun.ExcludeChangeLocations[i]) {
				t.Errorf("Expected excluded run change location to be '%v', but was '%v'", expectedRun.ExcludeChangeLocations[i], changeLocation)
			}
		}
	})
}
