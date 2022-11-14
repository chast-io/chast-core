package graph

import log "github.com/sirupsen/logrus"

type Node[T interface{}] struct {
	self         *T
	dependents   map[*Node[T]]bool
	dependencies map[*Node[T]]bool
}

func NewNode[T interface{}](run *T) *Node[T] {
	return &Node[T]{
		self:         run,
		dependents:   make(map[*Node[T]]bool),
		dependencies: make(map[*Node[T]]bool),
	}
}

func (n *Node[T]) addDependency(node *Node[T]) bool {
	if n.dependencies[node] {
		log.Warnf("dependency already exists: %v -> %v", n.self, node.self)

		return false
	}

	n.dependencies[node] = true
	node.dependents[n] = true

	return true
}

func (n *Node[T]) removeDependency(node *Node[T]) bool {
	if !n.dependencies[node] {
		return false
	}

	delete(n.dependencies, node)
	delete(node.dependents, n)

	return true
}
