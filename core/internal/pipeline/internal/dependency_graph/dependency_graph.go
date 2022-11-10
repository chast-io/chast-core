package dependencygraph

import (
	"chast.io/core/internal/run_model/pkg/model/refactoring"
	"github.com/pkg/errors"
)

var ErrCyclicDependency = errors.New("cyclic dependency detected")

func BuildExecutionOrder(runModel *refactoring.RunModel) ([][]*refactoring.Run, error) {
	executionOrder := make([][]*refactoring.Run, 0)

	coveredRunsCount := 0
	queue := buildDependencyGraph(runModel)

	for len(queue) > 0 {
		levelLen := len(queue)
		level := make([]*refactoring.Run, 0)

		for i := 0; i < levelLen; i++ {
			queueNode := queue[0]
			queue = queue[1:]

			level = append(level, queueNode.self)
			coveredRunsCount++

			for dependent := range queueNode.dependents {
				dependent.removeDependency(queueNode)

				if len(dependent.dependencies) == 0 {
					queue = append(queue, dependent)
				}
			}
		}

		executionOrder = append(executionOrder, level)
	}

	if coveredRunsCount != len(runModel.Run) {
		return nil, ErrCyclicDependency
	}

	return executionOrder, nil
}

func buildDependencyGraph(runModel *refactoring.RunModel) []*node {
	roots := make([]*node, 0)
	nodes := make([]*node, 0)
	nodesMap := make(map[*refactoring.Run]*node)

	for _, run := range runModel.Run {
		node := newNode(run)
		nodes = append(nodes, node)
		nodesMap[run] = node
	}

	for _, node := range nodes {
		if len(node.self.Dependencies) == 0 {
			roots = append(roots, node)
		} else {
			for _, dependency := range node.self.Dependencies {
				node.addDependency(nodesMap[dependency])
			}
		}
	}

	return roots
}
