package graph

type DoubleConnectedGraph[T interface{}] struct {
	nodes map[*Node[T]]bool
	roots map[*Node[T]]bool
}

func (graph *DoubleConnectedGraph[T]) AddNode(node *Node[T]) {
	graph.nodes[node] = true

	if len(node.dependencies) == 0 {
		graph.roots[node] = true
	}
}

func (graph *DoubleConnectedGraph[T]) RemoveNode(node *Node[T]) {
	delete(graph.nodes, node)
	delete(graph.roots, node)
}

func (graph *DoubleConnectedGraph[T]) AddEdge(node *Node[T], dependency *Node[T]) {
	node.addDependency(dependency)

	if graph.roots[node] {
		delete(graph.roots, node)
	}
}

func (graph *DoubleConnectedGraph[T]) RemoveEdge(node *Node[T], dependency *Node[T]) {
	node.removeDependency(dependency)

	if len(node.dependencies) == 0 {
		graph.roots[node] = true
	}
}

func (graph *DoubleConnectedGraph[T]) HasCycles() bool {
	visited := make(map[*Node[T]]bool)
	recStack := make(map[*Node[T]]bool)

	for rootNode, _ := range graph.roots {
		if graph.hasCyclesRecursive(rootNode, visited, recStack) {
			return true
		}
	}

	return false
}

func (graph *DoubleConnectedGraph[T]) hasCyclesRecursive(node *Node[T], visited map[*Node[T]]bool, recStack map[*Node[T]]bool) bool {
	if recStack[node] {
		return true
	}

	if visited[node] {
		return false
	}

	visited[node] = true
	recStack[node] = true

	for dependant := range node.dependents {
		if graph.hasCyclesRecursive(dependant, visited, recStack) {
			return true
		}
	}

	recStack[node] = false

	return false
}
