package graph_test

import (
	"testing"

	uut "chast.io/core/internal/internal_util/graph"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

// region: Helpers

func nodeDummyRun() *refactoring.Run {
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

	run := nodeDummyRun()

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

// region: AddDependency

//nolint:gocognit // nested test cases
func TestAddDependency(t *testing.T) {
	t.Parallel()

	t.Run("should add new dependency", func(t *testing.T) {
		t.Parallel()

		run := nodeDummyRun()
		node := uut.NewNode(run)

		dependencyRun := nodeDummyRun()
		dependencyNode := uut.NewNode(dependencyRun)

		response := node.AddDependency(dependencyNode)

		t.Run("should return true", func(t *testing.T) {
			t.Parallel()

			if response != true {
				t.Error("Expected response to be true, but was false")
			}
		})

		t.Run("should add dependency to Dependencies map", func(t *testing.T) {
			t.Parallel()

			if len(node.Dependencies) != 1 {
				t.Errorf("Expected Node Dependencies to contain 1 element, but was '%v'", node.Dependencies)
			}
			if node.Dependencies[dependencyNode] == false {
				t.Error("Expected Node Dependencies to contain dependencyNode, but did not")
			}
		})

		t.Run("should add dependent to Dependents map", func(t *testing.T) {
			t.Parallel()

			if len(dependencyNode.Dependents) != 1 {
				t.Errorf("Expected dependencyNode Dependents to have size 1, but was '%v'", dependencyNode.Dependents)
			}

			if dependencyNode.Dependents[node] == false {
				t.Error("Expected dependencyNode Dependents to contain Node, but did not")
			}
		})
	})

	t.Run("should not add existing dependency", func(t *testing.T) {
		t.Parallel()

		run := nodeDummyRun()
		node := uut.NewNode(run)

		dependencyRun := nodeDummyRun()
		dependencyNode := uut.NewNode[*refactoring.Run](dependencyRun)

		node.AddDependency(dependencyNode)

		response := node.AddDependency(dependencyNode)

		t.Run("should return false", func(t *testing.T) {
			t.Parallel()

			if response == true {
				t.Error("Expected Node.AddDependency to return false, but was true")
			}
		})

		t.Run("should not add dependency", func(t *testing.T) {
			t.Parallel()

			if len(node.Dependencies) != 1 {
				t.Errorf("Expected Node Dependencies to have size 1, but was '%v'", node.Dependencies)
			}
		})

		t.Run("should not add dependent", func(t *testing.T) {
			t.Parallel()

			if len(dependencyNode.Dependents) != 1 {
				t.Errorf("Expected dependencyNode Dependents to have size 1, but was '%v'", dependencyNode.Dependents)
			}
		})
	})
}

// endregion

// region: RemoveDependency

//nolint:gocognit // nested test
func TestRemoveDependency(t *testing.T) {
	t.Parallel()

	t.Run("should remove existing dependency", func(t *testing.T) {
		t.Parallel()

		// Prepare
		run := nodeDummyRun()
		node := uut.NewNode(run)

		dependencyRun := nodeDummyRun()
		dependencyNode := uut.NewNode(dependencyRun)

		node.AddDependency(dependencyNode)

		// Test
		response := node.RemoveDependency(dependencyNode)

		// Assert
		t.Run("should return true", func(t *testing.T) {
			t.Parallel()

			if response != true {
				t.Error("Expected response to be true, but was false")
			}
		})

		t.Run("should remove dependency from Dependencies map", func(t *testing.T) {
			t.Parallel()

			if len(node.Dependencies) != 0 {
				t.Errorf("Expected Node Dependencies to be empty, but was '%v'", node.Dependencies)
			}
		})

		t.Run("should remove dependent from Dependents map", func(t *testing.T) {
			t.Parallel()

			if len(dependencyNode.Dependents) != 0 {
				t.Errorf("Expected dependencyNode Dependents to be empty, but was '%v'", dependencyNode.Dependents)
			}
		})
	})

	t.Run("should not remove non-existing dependency", func(t *testing.T) {
		t.Parallel()

		t.Run("should do nothing if dependency does not exist", func(t *testing.T) {
			t.Parallel()

			// Prepare
			run := nodeDummyRun()
			node := uut.NewNode(run)

			dependencyRun := nodeDummyRun()
			dependencyNode := uut.NewNode(dependencyRun)

			otherDependencyRun := nodeDummyRun()
			otherDependencyNode := uut.NewNode(otherDependencyRun)

			otherDependencyNode.AddDependency(dependencyNode)

			node.AddDependency(otherDependencyNode)

			// Test
			response := node.RemoveDependency(dependencyNode)

			// Assert
			t.Run("should return false", func(t *testing.T) {
				t.Parallel()

				if response == true {
					t.Error("Expected response to be false, but was true")
				}
			})

			t.Run("should not remove dependency", func(t *testing.T) {
				t.Parallel()

				if len(node.Dependencies) != 1 {
					t.Errorf("Expected Node Dependencies to be empty, but was '%v'", node.Dependencies)
				}

				if node.Dependencies[otherDependencyNode] == false {
					t.Error("Expected Node Dependencies to contain otherDependencyNode, but did not")
				}
			})

			t.Run("should not remove dependent", func(t *testing.T) {
				t.Parallel()

				if len(dependencyNode.Dependents) != 1 {
					t.Errorf("Expected dependencyNode Dependents to be 1, but was '%v'", dependencyNode.Dependents)
				}

				if otherDependencyNode.Dependencies[dependencyNode] == false {
					t.Error("Expected otherDependencyNode Dependencies to contain Node, but did not")
				}
			})
		})
	})

	t.Run("should remove dependency if dependency exists", func(t *testing.T) {
		t.Parallel()

		// Prepare
		run := nodeDummyRun()
		node := uut.NewNode(run)

		dependencyRun := nodeDummyRun()
		dependencyNode := uut.NewNode(dependencyRun)

		node.AddDependency(dependencyNode)

		// Test
		response := node.RemoveDependency(dependencyNode)

		t.Run("should return true", func(t *testing.T) {
			t.Parallel()

			if response != true {
				t.Error("Expected response to be true, but was false")
			}
		})

		t.Run("should remove dependency from Dependencies map", func(t *testing.T) {
			t.Parallel()

			if len(node.Dependencies) != 0 {
				t.Errorf("Expected Node Dependencies to be empty, but was '%v'", node.Dependencies)
			}
		})

		t.Run("should remove dependent from Dependents map", func(t *testing.T) {
			t.Parallel()

			if len(dependencyNode.Dependents) != 0 {
				t.Errorf("Expected dependencyNode Dependents to be empty, but was '%v'", dependencyNode.Dependents)
			}
		})
	})
}

// endregion
