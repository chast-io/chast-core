package refactroingdependencygraph

import (
	"chast.io/core/internal/internal_util/graph"
	recipemodel "chast.io/core/internal/recipe/pkg/model"
)

func BuildDependencyGraph(runs []recipemodel.Run) *graph.DoubleConnectedGraph[*recipemodel.Run] {
	nodesMap := make(map[string]*graph.Node[*recipemodel.Run])

	runGraph := graph.NewDoubleConnectedGraph[*recipemodel.Run]()

	for i, run := range runs {
		node := graph.NewNode(&runs[i])
		runGraph.AddNode(node)
		nodesMap[run.ID] = node
	}

	for node := range runGraph.Nodes() {
		for _, dependency := range node.Self.Dependencies {
			dependencyNode := nodesMap[dependency]
			if dependencyNode == nil {
				continue // this can happen if the dependency is a run that is not part of the run model due to a filter
			}

			runGraph.AddEdge(node, dependencyNode)
		}
	}

	return runGraph
}
