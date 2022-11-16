package graph

// TODO make order deterministic
type DoubleConnectedGraph[T interface{}] struct {
	Nodes map[*Node[T]]bool
	Roots map[*Node[T]]bool
}

func NewDoubleConnectedGraph[T interface{}]() *DoubleConnectedGraph[T] {
	return &DoubleConnectedGraph[T]{
		Nodes: make(map[*Node[T]]bool),
		Roots: make(map[*Node[T]]bool),
	}
}

func (graph *DoubleConnectedGraph[T]) AddNode(node *Node[T]) {
	graph.Nodes[node] = true

	if len(node.Dependencies) == 0 {
		graph.Roots[node] = true
	}
}

func (graph *DoubleConnectedGraph[T]) RemoveNode(node *Node[T]) {
	delete(graph.Nodes, node)
	delete(graph.Roots, node)
}

func (graph *DoubleConnectedGraph[T]) AddEdge(node *Node[T], dependency *Node[T]) bool {
	graph.AddNode(node)
	graph.AddNode(dependency)

	success := node.AddDependency(dependency)

	if success && graph.Roots[node] {
		delete(graph.Roots, node)
	}

	return success
}

func (graph *DoubleConnectedGraph[T]) RemoveEdge(node *Node[T], dependency *Node[T]) bool {
	success := node.RemoveDependency(dependency)

	if success && len(node.Dependencies) == 0 {
		graph.Roots[node] = true
	}

	return success
}

func (graph *DoubleConnectedGraph[T]) HasCycles() bool {
	visited := make(map[*Node[T]]bool)
	recStack := make(map[*Node[T]]bool)

	for rootNode, _ := range graph.Roots {
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

	for dependant := range node.Dependents {
		if graph.hasCyclesRecursive(dependant, visited, recStack) {
			return true
		}
	}

	recStack[node] = false

	return false
}
