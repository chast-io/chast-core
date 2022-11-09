package dependencygraph

import (
	"log"

	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

func BuildExecutionOrder(runModel *refactoring.RunModel) [][]*refactoring.Run {
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
				dependent.RemoveDependency(queueNode)

				if len(dependent.dependencies) == 0 {
					queue = append(queue, dependent)
				}
			}
		}

		executionOrder = append(executionOrder, level)
	}

	if coveredRunsCount != len(runModel.Run) {
		log.Fatal("cyclic dependency detected")
	}

	return executionOrder
}

func buildDependencyGraph(runModel *refactoring.RunModel) []*Node {
	roots := make([]*Node, 0)
	nodes := make([]*Node, 0)
	nodesMap := make(map[*refactoring.Run]*Node)

	for _, run := range runModel.Run {
		node := NewNode(run)
		nodes = append(nodes, node)
		nodesMap[run] = node
	}

	for _, node := range nodes {
		if len(node.self.Dependencies) == 0 {
			roots = append(roots, node)
		} else {
			for _, dependency := range node.self.Dependencies {
				node.AddDependency(nodesMap[dependency])
			}
		}
	}

	return roots
}

func nodesToRuns(nodes []*Node) []*refactoring.Run {
	runs := make([]*refactoring.Run, 0)
	for _, node := range nodes {
		runs = append(runs, node.self)
	}
	return runs
}
