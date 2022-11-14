package dependencygraph

import (
	"chast.io/core/internal/internal_util/graph"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
	"github.com/pkg/errors"
)

var ErrCyclicDependency = errors.New("cyclic dependency detected")

func BuildExecutionOrder(runModel *refactoring.RunModel) ([][]*refactoring.Run, error) {
	executionOrder := make([][]*refactoring.Run, 0)

	dependencyGraph := buildDependencyGraph(runModel)

	if dependencyGraph.hasCycles() {
		return nil, ErrCyclicDependency
	}

	queue := dependencyGraph

	for len(queue) > 0 {
		levelLen := len(queue)
		level := make([]*refactoring.Run, 0)

		for i := 0; i < levelLen; i++ {
			queueNode := queue[0]
			queue = queue[1:]

			level = append(level, queueNode.self)

			for dependent := range queueNode.dependents {
				dependent.removeDependency(queueNode)

				if len(dependent.dependencies) == 0 {
					queue = append(queue, dependent)
				}
			}
		}

		executionOrder = append(executionOrder, level)
	}

	return executionOrder, nil
}

func buildDependencyGraph(runModel *refactoring.RunModel) *graph.DoubleConnectedGraph[refactoring.Run] {
	roots := make([]*graph.node, 0)
	nodes := make([]*graph.node, 0)
	nodesMap := make(map[*refactoring.Run]*graph.node)

	for _, run := range runModel.Run {
		node := graph.NewNode(run)
		nodes = append(nodes, node)
		nodesMap[run] = node
	}

	for _, node := range nodes {
		if len(node.self.Dependencies) == 0 {
			roots = append(roots, node)
		} else {
			for _, dependency := range node.self.Dependencies {
				dependencyNode := nodesMap[dependency]
				if dependencyNode == nil {
					continue // this can happen if the dependency is a run that is not part of the run model due to a filter
				}
				node.addDependency(dependencyNode)
			}
		}
	}

	return roots
}
