package graph_test

import (
	"testing"

	uut "chast.io/core/internal/internal_util/graph"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

// region: Helpers
func graphDummyGraph() *uut.DoubleConnectedGraph[*refactoring.Run] {
	graph := uut.NewDoubleConnectedGraph[*refactoring.Run]()

	return graph
}

func graphDummyRunNoDependencies() *refactoring.Run {
	return &refactoring.Run{
		ID:                 "graphDummyRunNoDependencies",
		Dependencies:       make([]*refactoring.Run, 0),
		SupportedLanguages: []string{"java"},
		Docker:             &refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              &refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            &refactoring.Command{}, //nolint:exhaustruct // not required for test
	}
}

func graphDummyRunWithDependencies() *refactoring.Run {
	return &refactoring.Run{
		ID:                 "graphDummyRunWithDependencies",
		Dependencies:       []*refactoring.Run{graphDummyRunNoDependencies()},
		SupportedLanguages: []string{"java"},
		Docker:             &refactoring.Docker{},  //nolint:exhaustruct // not required for test
		Local:              &refactoring.Local{},   //nolint:exhaustruct // not required for test
		Command:            &refactoring.Command{}, //nolint:exhaustruct // not required for test
	}
}

// endregion

// region: AddNode

//nolint:gocognit // nested test method
func TestAddNode(t *testing.T) {
	t.Parallel()

	t.Run("No dependencies", func(t *testing.T) {
		t.Parallel()

		graph := graphDummyGraph()
		node := uut.NewNode[*refactoring.Run](graphDummyRunNoDependencies())

		graph.AddNode(node)

		t.Run("should add node to Nodes", func(t *testing.T) {
			nodes := graph.Nodes()
			retrievedNode := <-nodes
			_, continuing := <-nodes

			if retrievedNode == nil {
				t.Fatalf("Expected node to be added to nodes, but was nil")
			}

			if continuing {
				t.Fatalf("Expected node to be only node in nodes, but was not")
			}

			if retrievedNode != node {
				t.Fatalf("Expected node to be added to nodes, but was '%v'", retrievedNode)
			}
		})

		t.Run("should add node to Roots", func(t *testing.T) {
			nodes := graph.Nodes()
			retrievedNode := <-nodes
			_, continuing := <-nodes

			if retrievedNode == nil {
				t.Fatalf("Expected node to be added to roots, but was nil")
			}

			if continuing {
				t.Fatalf("Expected node to be only node in roots, but was not")
			}

			if retrievedNode != node {
				t.Fatalf("Expected node to be added to roots, but was '%v'", retrievedNode)
			}
		})
	})

	t.Run("With dependencies", func(t *testing.T) {
		t.Parallel()

		graph := graphDummyGraph()
		node := uut.NewNode[*refactoring.Run](graphDummyRunWithDependencies())
		nodeDependency := uut.NewNode[*refactoring.Run](graphDummyRunNoDependencies())

		node.AddDependency(nodeDependency)

		graph.AddNode(node)

		t.Run("should add node and nodeDependency to Nodes", func(t *testing.T) {
			t.Parallel()

			nodes := graph.Nodes()
			retrievedNode := <-nodes

			if retrievedNode == nil {
				t.Fatalf("Expected node to be added to nodes, but was nil")
			}

			if retrievedNode != node {
				t.Fatalf("Expected node to be added to nodes first, but was '%v'", retrievedNode)
			}

			retrievedNode, continuing := <-nodes

			if !continuing {
				t.Fatalf("Expected node dependency to be second node in nodes, but ended early")
			}

			if retrievedNode == nil {
				t.Fatalf("Expected node dependency to be added to nodes, but was nil")
			}

			if retrievedNode != nodeDependency {
				t.Fatalf("Expected node dependency to be added to nodes second, but was '%v'", retrievedNode.Self.ID)
			}

			_, continuing = <-nodes

			if continuing {
				t.Fatalf("Expected to contain only two nodes in nodes, but was more")
			}
		})

		t.Run("should not add node, but nodeDependency to Roots", func(t *testing.T) {
			t.Parallel()

			nodes := graph.Roots()
			retrievedNode := <-nodes
			_, continuing := <-nodes

			if retrievedNode != nodeDependency {
				t.Errorf("Expected node to not be added to roots, but was '%v'", retrievedNode)
			}

			if continuing {
				t.Error("Expected roots to be empty, but was not")
			}
		})
	})
}

// endregion

// region: RemoveNode

//nolint:gocognit // nested test method
func TestRemoveNode(t *testing.T) {
	t.Parallel()

	t.Run("No dependencies", func(t *testing.T) {
		t.Parallel()

		graph := graphDummyGraph()
		node := uut.NewNode[*refactoring.Run](graphDummyRunNoDependencies())

		graph.AddNode(node)

		graph.RemoveNode(node)

		t.Run("should remove node from Nodes", func(t *testing.T) {
			t.Parallel()

			nodes := graph.Nodes()
			_, continuing := <-nodes

			if continuing {
				t.Fatalf("Expected nodes to be empty, but was not")
			}
		})

		t.Run("should remove node from Roots", func(t *testing.T) {
			t.Parallel()

			nodes := graph.Roots()
			_, continuing := <-nodes

			if continuing {
				t.Fatalf("Expected roots to be empty, but was not")
			}
		})
	})

	t.Run("With dependencies", func(t *testing.T) {
		t.Parallel()

		graph := graphDummyGraph()
		node := uut.NewNode[*refactoring.Run](graphDummyRunWithDependencies())
		nodeDependency := uut.NewNode[*refactoring.Run](graphDummyRunNoDependencies())

		node.AddDependency(nodeDependency)

		graph.AddNode(node)

		graph.RemoveNode(node)

		t.Run("should remove node from Nodes", func(t *testing.T) {
			t.Parallel()

			nodes := graph.Nodes()
			presentNode, continuing := <-nodes

			if !continuing {
				t.Fatalf("Expected nodes to contain node dependency, but was empty")
			}

			if presentNode != nodeDependency {
				t.Fatalf("Expected nodes to contain node dependency, but was '%v'", presentNode)
			}

			_, continuing = <-nodes

			if continuing {
				t.Fatalf("Expected nodes to be empty, but was not")
			}
		})

		t.Run("should remove node from nodeDependency dependencies", func(t *testing.T) {
			t.Parallel()

			if len(nodeDependency.Dependencies) != 0 {
				t.Fatalf("Expected node dependency to not have dependencies, but was '%v'", nodeDependency.Dependencies)
			}

			if len(nodeDependency.Dependents) != 0 {
				t.Fatalf("Expected node dependency to not have dependents, but was '%v'", nodeDependency.Dependents)
			}
		})
	})

	t.Run("With dependents", func(t *testing.T) {
		t.Parallel()

		graph := graphDummyGraph()
		node := uut.NewNode[*refactoring.Run](graphDummyRunWithDependencies())
		nodeDependency := uut.NewNode[*refactoring.Run](graphDummyRunNoDependencies())

		node.AddDependency(nodeDependency)

		graph.AddNode(node)

		graph.RemoveNode(nodeDependency)

		t.Run("should remove node from Nodes", func(t *testing.T) {
			t.Parallel()

			nodes := graph.Nodes()
			presentNode, continuing := <-nodes

			if !continuing {
				t.Fatalf("Expected nodes to contain node, but was empty")
			}

			if presentNode != node {
				t.Fatalf("Expected nodes to contain node, but was '%v'", presentNode)
			}

			_, continuing = <-nodes

			if continuing {
				t.Fatalf("Expected nodes to be empty, but was not")
			}
		})

		t.Run("should remove node from node dependencies", func(t *testing.T) {
			t.Parallel()

			if len(node.Dependencies) != 0 {
				t.Fatalf("Expected node to not have dependencies, but was '%v'", node.Dependencies)
			}

			if len(node.Dependents) != 0 {
				t.Fatalf("Expected node to not have dependents, but was '%v'", node.Dependents)
			}
		})
	})
}

// endregion

// region: AddEdge

func TestAddEdge(t *testing.T) {
	t.Parallel()

	t.Run("should remove source node from Roots", func(t *testing.T) {
		t.Parallel()

		graph := graphDummyGraph()
		node := uut.NewNode[*refactoring.Run](graphDummyRunNoDependencies())
		nodeDependency := uut.NewNode[*refactoring.Run](graphDummyRunNoDependencies())

		graph.AddNode(node)
		graph.AddNode(nodeDependency)

		roots := graph.Roots()
		nodes1, continuing := <-roots

		if nodes1 == nil || !continuing {
			t.Fatalf("Expected roots to contain node dependency, but was empty")
		}

		nodes2, continuing := <-roots

		if nodes2 == nil || !continuing {
			t.Fatalf("Expected roots to contain node dependency, but was empty")
		}

		graph.AddEdge(node, nodeDependency)

		roots = graph.Roots()
		presentNode := <-roots
		_, continuing = <-roots

		if presentNode != nodeDependency {
			t.Fatalf("Expected roots to contain node dependency, but was empty")
		}

		if continuing {
			t.Fatalf("Expected roots to be empty, but was not")
		}
	})
}

// endregion

// region: RemoveEdge

//nolint:gocognit // nested test
func TestRemoveEdge(t *testing.T) {
	t.Parallel()

	t.Run("should remove existing dependency", func(t *testing.T) {
		t.Parallel()

		// Prepare
		run := graphDummyRunNoDependencies()
		node := uut.NewNode(run)

		dependencyRun := graphDummyRunNoDependencies()
		dependencyNode := uut.NewNode(dependencyRun)
		graph := graphDummyGraph()
		graph.AddEdge(node, dependencyNode)

		// Test
		response := graph.RemoveEdge(node, dependencyNode)

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
			run := graphDummyRunNoDependencies()
			node := uut.NewNode(run)

			graph := graphDummyGraph()

			dependencyRun := graphDummyRunNoDependencies()
			dependencyNode := uut.NewNode(dependencyRun)

			otherDependencyRun := graphDummyRunNoDependencies()
			otherDependencyNode := uut.NewNode(otherDependencyRun)

			graph.AddEdge(otherDependencyNode, dependencyNode)

			graph.AddEdge(node, otherDependencyNode)

			// Test
			response := graph.RemoveEdge(node, dependencyNode)

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
		run := graphDummyRunNoDependencies()
		node := uut.NewNode(run)

		dependencyRun := graphDummyRunNoDependencies()
		dependencyNode := uut.NewNode(dependencyRun)

		graph := graphDummyGraph()
		graph.AddEdge(node, dependencyNode)

		// Test
		response := graph.RemoveEdge(node, dependencyNode)

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

// region: HasCycles

func TestHasCycles(t *testing.T) {
	t.Parallel()

	t.Run("should return false if no cycles exist", func(t *testing.T) {
		t.Parallel()

		// Prepare
		run := graphDummyRunNoDependencies()
		node := uut.NewNode(run)

		graph := graphDummyGraph()
		graph.AddNode(node)

		// Test
		response := graph.HasCycles()

		// Assert
		if response != false {
			t.Error("Expected response to be false, but was true")
		}
	})

	t.Run("should return true if cycle exists at start", func(t *testing.T) {
		t.Parallel()

		// Prepare
		run := graphDummyRunNoDependencies()
		node := uut.NewNode(run)

		dependencyRun := graphDummyRunNoDependencies()
		dependencyNode := uut.NewNode(dependencyRun)

		graph := graphDummyGraph()
		graph.AddEdge(node, dependencyNode)
		graph.AddEdge(dependencyNode, node)

		// Test
		response := graph.HasCycles()

		// Assert
		if response != true {
			t.Error("Expected response to be true, but was false")
		}
	})

	t.Run("should return true if cycle exists in multiple nodes at start", func(t *testing.T) {
		t.Parallel()

		// Prepare
		run := graphDummyRunNoDependencies()
		node := uut.NewNode(run)

		dependencyRun := graphDummyRunNoDependencies()
		dependencyNode := uut.NewNode(dependencyRun)

		otherDependencyRun := graphDummyRunNoDependencies()
		otherDependencyNode := uut.NewNode(otherDependencyRun)

		graph := graphDummyGraph()
		graph.AddEdge(node, dependencyNode)
		graph.AddEdge(dependencyNode, otherDependencyNode)
		graph.AddEdge(otherDependencyNode, node)

		// Test
		response := graph.HasCycles()

		// Assert
		if response != true {
			t.Error("Expected response to be true, but was false")
		}
	})

	t.Run("should return true if cycle exists in the middle", func(t *testing.T) {
		t.Parallel()

		// Prepare
		run1 := graphDummyRunNoDependencies()
		node1 := uut.NewNode(run1)

		run2 := graphDummyRunNoDependencies()
		node2 := uut.NewNode(run2)

		run3 := graphDummyRunNoDependencies()
		node3 := uut.NewNode(run3)

		run4 := graphDummyRunNoDependencies()
		node4 := uut.NewNode(run4)

		graph := graphDummyGraph()
		graph.AddNode(node1)
		graph.AddNode(node2)
		graph.AddNode(node3)
		graph.AddNode(node4)

		graph.AddEdge(node2, node1)
		graph.AddEdge(node2, node4)
		graph.AddEdge(node3, node2)
		graph.AddEdge(node4, node3)

		// Test
		response := graph.HasCycles()

		// Assert
		if response != true {
			t.Error("Expected response to be true, but was false")
		}
	})

	t.Run("should return false if no cycle but multiple paths exists", func(t *testing.T) {
		t.Parallel()

		// Prepare
		run1 := graphDummyRunNoDependencies()
		node1 := uut.NewNode(run1)

		run2 := graphDummyRunNoDependencies()
		node2 := uut.NewNode(run2)

		run3 := graphDummyRunNoDependencies()
		node3 := uut.NewNode(run3)

		run4 := graphDummyRunNoDependencies()
		node4 := uut.NewNode(run4)

		graph := graphDummyGraph()
		graph.AddNode(node1)
		graph.AddNode(node2)
		graph.AddNode(node3)
		graph.AddNode(node4)

		graph.AddEdge(node2, node1)
		graph.AddEdge(node3, node2)
		graph.AddEdge(node4, node3)
		graph.AddEdge(node4, node2)

		// Test
		response := graph.HasCycles()

		// Assert
		if response != false {
			t.Error("Expected response to be false, but was true")
		}
	})
}

// endregion
