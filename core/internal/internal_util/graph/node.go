package graph

import log "github.com/sirupsen/logrus"

type Node[T interface{}] struct {
	Self         T
	Dependents   map[*Node[T]]bool
	Dependencies map[*Node[T]]bool
}

func NewNode[T interface{}](run T) *Node[T] {
	return &Node[T]{
		Self:         run,
		Dependents:   make(map[*Node[T]]bool),
		Dependencies: make(map[*Node[T]]bool),
	}
}

func (n *Node[T]) AddDependency(node *Node[T]) bool {
	if n.Dependencies[node] {
		log.Warnf("dependency already exists: %v -> %v", n.Self, node.Self)

		return false
	}

	n.Dependencies[node] = true
	node.Dependents[n] = true

	return true
}

func (n *Node[T]) RemoveDependency(node *Node[T]) bool {
	if !n.Dependencies[node] {
		return false
	}

	delete(n.Dependencies, node)
	delete(node.Dependents, n)

	return true
}
