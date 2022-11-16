package graph_test

import (
	uut "chast.io/core/internal/internal_util/graph"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
	"testing"
)

// region: Helpers

func NodeDummyRun() *refactoring.Run {
	return &refactoring.Run{
		ID:                 "runId",
		Dependencies:       make([]*refactoring.Run, 0),
		SupportedLanguages: []string{"java"},
		Docker:             &refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              &refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            &refactoring.Command{}, //nolint:exhaustruct // not required for test
	}
}

// endregion

// region: NewNode
func TestNewNode(t *testing.T) {
	t.Parallel()

	run := NodeDummyRun()

	uutNode := uut.NewNode[*refactoring.Run](run)

	t.Run("should set Self", func(t *testing.T) {
		t.Parallel()
		if uutNode.Self != run {
			t.Errorf("Expected Node Self to be '%v', but was '%v'", run, uutNode.Self)
		}
	})

	t.Run("should initialize Dependents map", func(t *testing.T) {
		t.Parallel()
		if uutNode.Dependents == nil {
			t.Error("Expected Node Dependents to be set, but was nil")
		}
		if len(uutNode.Dependents) != 0 {
			t.Errorf("Expected Node Dependents to be empty, but was '%v'", uutNode.Dependents)
		}
	})

	t.Run("should initialize Dependencies map", func(t *testing.T) {
		t.Parallel()
		if uutNode.Dependencies == nil {
			t.Error("Expected Node Dependencies to be set, but was nil")
		}
		if len(uutNode.Dependencies) != 0 {
			t.Errorf("Expected Node Dependencies to be empty, but was '%v'", uutNode.Dependencies)
		}
	})
}

// endregion
