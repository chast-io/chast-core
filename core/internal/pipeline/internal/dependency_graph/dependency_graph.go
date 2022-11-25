package dependencygraph

import (
	"chast.io/core/internal/internal_util/graph"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
	"github.com/joomcode/errorx"
)

func BuildExecutionOrder(runModel *refactoring.RunModel) ([][]*refactoring.Run, error) {
	executionOrder := make([][]*refactoring.Run, 0)

	dependencyGraph := buildDependencyGraph(runModel)

	if dependencyGraph.HasCycles() {
		return nil, errorx.InternalError.New("Cyclic dependency detected")
	}

	queue := make([]*graph.Node[*refactoring.Run], 0)
	for node := range dependencyGraph.Roots() {
		queue = append(queue, node)
	}

	for len(queue) > 0 {
		levelLen := len(queue)
		level := make([]*refactoring.Run, 0)

		for i := 0; i < levelLen; i++ {
			queueNode := queue[0]
			queue = queue[1:]

			level = append(level, queueNode.Self)

			for dependent := range queueNode.Dependents {
				dependent.RemoveDependency(queueNode)

				if len(dependent.Dependencies) == 0 {
					queue = append(queue, dependent)
				}
			}
		}

		executionOrder = append(executionOrder, level)
	}

	return executionOrder, nil
}

func buildDependencyGraph(runModel *refactoring.RunModel) *graph.DoubleConnectedGraph[*refactoring.Run] {
	nodesMap := make(map[*refactoring.Run]*graph.Node[*refactoring.Run])

	runGraph := graph.NewDoubleConnectedGraph[*refactoring.Run]()

	for _, run := range runModel.Run {
		node := graph.NewNode[*refactoring.Run](run)
		runGraph.AddNode(node)
		nodesMap[run] = node
	}

	for node := range runGraph.Nodes() {
		for _, dependency := range node.Self.Dependencies {
			dependencyNode := nodesMap[dependency]
			if dependencyNode == nil {
				runGraph.RemoveNode(node)
				// this can happen if the dependency is a run that is not part of the run model due to a filter
				break
			}

			runGraph.AddEdge(node, dependencyNode)
		}
	}

	return runGraph
}
