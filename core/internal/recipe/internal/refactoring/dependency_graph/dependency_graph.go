package refactroingdependencygraph

import (
	recipemodel "chast.io/core/internal/recipe/pkg/model"
	"chast.io/core/internal/run_model/pkg/model/refactoring"
)

func buildDependencyGraph(runModel *recipemodel.RefactoringRecipe) []*node {
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

func hasCycles(nodes []*node) bool {
	visited := make(map[*node]bool)
	recStack := make(map[*node]bool)

	for _, rootNode := range nodes {
		if hasCyclesRecursive(rootNode, visited, recStack) {
			return true
		}
	}

	return false
}

func hasCyclesRecursive(node *node, visited map[*node]bool, recStack map[*node]bool) bool {
	if recStack[node] {
		return true
	}

	if visited[node] {
		return false
	}

	visited[node] = true
	recStack[node] = true

	for dependant := range node.dependents {
		if hasCyclesRecursive(dependant, visited, recStack) {
			return true
		}
	}

	recStack[node] = false

	return false
}
