package graph

import orderedmap "github.com/wk8/go-ordered-map/v2"

type DoubleConnectedGraph[T interface{}] struct {
	nodes *orderedmap.OrderedMap[*Node[T], bool]
	roots *orderedmap.OrderedMap[*Node[T], bool]
}

func NewDoubleConnectedGraph[T interface{}]() *DoubleConnectedGraph[T] {
	return &DoubleConnectedGraph[T]{
		nodes: orderedmap.New[*Node[T], bool](),
		roots: orderedmap.New[*Node[T], bool](),
	}
}

func (graph *DoubleConnectedGraph[T]) Roots() chan *Node[T] {
	channel := make(chan *Node[T])

	go func() {
		for pair := graph.roots.Oldest(); pair != nil; pair = pair.Next() {
			channel <- pair.Key
		}
		close(channel)
	}()

	return channel
}

func (graph *DoubleConnectedGraph[T]) Nodes() chan *Node[T] {
	channel := make(chan *Node[T])

	go func() {
		for pair := graph.nodes.Oldest(); pair != nil; pair = pair.Next() {
			channel <- pair.Key
		}
		close(channel)
	}()

	return channel
}

func (graph *DoubleConnectedGraph[T]) AddNode(node *Node[T]) {
	if _, present := graph.nodes.Get(node); present {
		return
	}

	graph.nodes.Set(node, true)

	if len(node.Dependencies) == 0 {
		graph.roots.Set(node, true)
	} else {
		for dependency := range node.Dependencies {
			graph.AddNode(dependency)
			graph.addEdgeToExistingNodes(node, dependency)
		}
	}
}

func (graph *DoubleConnectedGraph[T]) RemoveNode(node *Node[T]) {
	graph.nodes.Delete(node)
	graph.roots.Delete(node)

	for dependency := range node.Dependencies {
		node.RemoveDependency(dependency)
	}

	for dependent := range node.Dependents {
		dependent.RemoveDependency(node)
	}
}

func (graph *DoubleConnectedGraph[T]) AddEdge(node *Node[T], dependency *Node[T]) bool {
	graph.AddNode(node)
	graph.AddNode(dependency)

	return graph.addEdgeToExistingNodes(node, dependency)
}

func (graph *DoubleConnectedGraph[T]) addEdgeToExistingNodes(node *Node[T], dependency *Node[T]) bool {
	success := node.AddDependency(dependency)

	if success {
		graph.roots.Delete(node)
	}

	return success
}

func (graph *DoubleConnectedGraph[T]) RemoveEdge(node *Node[T], dependency *Node[T]) bool {
	success := node.RemoveDependency(dependency)

	if success && len(node.Dependencies) == 0 {
		graph.roots.Set(node, true)
	}

	return success
}

func (graph *DoubleConnectedGraph[T]) HasCycles() bool {
	visited := make(map[*Node[T]]bool)
	recStack := make(map[*Node[T]]bool)

	if graph.roots.Len() == 0 {
		return true
	}

	for pair := graph.roots.Oldest(); pair != nil; pair = pair.Next() {
		if graph.hasCyclesRecursive(pair.Key, visited, recStack) {
			return true
		}
	}

	return false
}

func (graph *DoubleConnectedGraph[T]) hasCyclesRecursive(
	node *Node[T],
	visited map[*Node[T]]bool,
	recStack map[*Node[T]]bool,
) bool {
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
