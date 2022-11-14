package graph

import (
	"testing"

	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

// region: Helpers
func GraphDummyRun() *refactoring.Run {
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

	run := GraphDummyRun()

	uutNode := NewNode(run)

	t.Run("should set self", func(t *testing.T) {
		t.Parallel()
		if uutNode.self != run {
			t.Errorf("Expected Node self to be '%v', but was '%v'", run, uutNode.self)
		}
	})

	t.Run("should initialize dependents map", func(t *testing.T) {
		t.Parallel()
		if uutNode.dependents == nil {
			t.Error("Expected Node dependents to be set, but was nil")
		}
		if len(uutNode.dependents) != 0 {
			t.Errorf("Expected Node dependents to be empty, but was '%v'", uutNode.dependents)
		}
	})

	t.Run("should initialize dependencies map", func(t *testing.T) {
		t.Parallel()
		if uutNode.dependencies == nil {
			t.Error("Expected Node dependencies to be set, but was nil")
		}
		if len(uutNode.dependencies) != 0 {
			t.Errorf("Expected Node dependencies to be empty, but was '%v'", uutNode.dependencies)
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

		run := GraphDummyRun()
		node := NewNode(run)

		dependencyRun := GraphDummyRun()
		dependencyNode := NewNode(dependencyRun)

		response := node.addDependency(dependencyNode)

		t.Run("should return true", func(t *testing.T) {
			t.Parallel()

			if response != true {
				t.Error("Expected response to be true, but was false")
			}
		})

		t.Run("should add dependency to dependencies map", func(t *testing.T) {
			t.Parallel()

			if len(node.dependencies) != 1 {
				t.Errorf("Expected Node dependencies to contain 1 element, but was '%v'", node.dependencies)
			}
			if node.dependencies[dependencyNode] == false {
				t.Error("Expected Node dependencies to contain dependencyNode, but did not")
			}
		})

		t.Run("should add dependent to dependents map", func(t *testing.T) {
			t.Parallel()

			if len(dependencyNode.dependents) != 1 {
				t.Errorf("Expected dependencyNode dependents to have size 1, but was '%v'", dependencyNode.dependents)
			}

			if dependencyNode.dependents[node] == false {
				t.Error("Expected dependencyNode dependents to contain Node, but did not")
			}
		})
	})

	t.Run("should not add existing dependency", func(t *testing.T) {
		t.Parallel()

		run := GraphDummyRun()
		node := NewNode(run)

		dependencyRun := GraphDummyRun()
		dependencyNode := NewNode(dependencyRun)
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

			if len(node.dependencies) != 1 {
				t.Errorf("Expected Node dependencies to have size 1, but was '%v'", node.dependencies)
			}
		})

		t.Run("should not add dependent", func(t *testing.T) {
			t.Parallel()

			if len(dependencyNode.dependents) != 1 {
				t.Errorf("Expected dependencyNode dependents to have size 1, but was '%v'", dependencyNode.dependents)
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
		run := GraphDummyRun()
		node := NewNode(run)

		dependencyRun := GraphDummyRun()
		dependencyNode := NewNode(dependencyRun)
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

		t.Run("should remove dependency from dependencies map", func(t *testing.T) {
			t.Parallel()

			if len(node.dependencies) != 0 {
				t.Errorf("Expected Node dependencies to be empty, but was '%v'", node.dependencies)
			}
		})

		t.Run("should remove dependent from dependents map", func(t *testing.T) {
			t.Parallel()

			if len(dependencyNode.dependents) != 0 {
				t.Errorf("Expected dependencyNode dependents to be empty, but was '%v'", dependencyNode.dependents)
			}
		})
	})

	t.Run("should not remove non-existing dependency", func(t *testing.T) {
		t.Parallel()

		t.Run("should do nothing if dependency does not exist", func(t *testing.T) {
			t.Parallel()

			// Prepare
			run := GraphDummyRun()
			node := NewNode(run)

			dependencyRun := GraphDummyRun()
			dependencyNode := NewNode(dependencyRun)

			otherDependencyRun := GraphDummyRun()
			otherDependencyNode := NewNode(otherDependencyRun)

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

				if len(node.dependencies) != 1 {
					t.Errorf("Expected Node dependencies to be empty, but was '%v'", node.dependencies)
				}

				if node.dependencies[otherDependencyNode] == false {
					t.Error("Expected Node dependencies to contain otherDependencyNode, but did not")
				}
			})

			t.Run("should not remove dependent", func(t *testing.T) {
				t.Parallel()

				if len(dependencyNode.dependents) != 1 {
					t.Errorf("Expected dependencyNode dependents to be 1, but was '%v'", dependencyNode.dependents)
				}

				if otherDependencyNode.dependencies[dependencyNode] == false {
					t.Error("Expected otherDependencyNode dependencies to contain Node, but did not")
				}
			})
		})
	})

	t.Run("should remove dependency if dependency exists", func(t *testing.T) {
		t.Parallel()

		// Prepare
		run := GraphDummyRun()
		node := NewNode(run)

		dependencyRun := GraphDummyRun()
		dependencyNode := NewNode(dependencyRun)

		node.AddDependency(dependencyNode)

		// Test
		response := node.RemoveDependency(dependencyNode)

		t.Run("should return true", func(t *testing.T) {
			t.Parallel()

			if response != true {
				t.Error("Expected response to be true, but was false")
			}
		})

		t.Run("should remove dependency from dependencies map", func(t *testing.T) {
			t.Parallel()

			if len(node.dependencies) != 0 {
				t.Errorf("Expected Node dependencies to be empty, but was '%v'", node.dependencies)
			}
		})

		t.Run("should remove dependent from dependents map", func(t *testing.T) {
			t.Parallel()

			if len(dependencyNode.dependents) != 0 {
				t.Errorf("Expected dependencyNode dependents to be empty, but was '%v'", dependencyNode.dependents)
			}
		})
	})
}

// endregion
