package dependencygraph

import (
	"chast.io/core/internal/run_model/pkg/model/refactoring"
	log "github.com/sirupsen/logrus"
)

type Node struct {
	self         *refactoring.Run
	dependents   map[*Node]bool
	dependencies map[*Node]bool
}

func NewNode(run *refactoring.Run) *Node {
	return &Node{
		self:         run,
		dependents:   make(map[*Node]bool),
		dependencies: make(map[*Node]bool),
	}
}

func (n *Node) AddDependency(node *Node) bool {
	if n.dependencies[node] {
		log.Warnf("dependency already exists: %v -> %v", n.self.ID, node.self.ID)

		return false
	}

	n.dependencies[node] = true
	node.dependents[n] = true

	return true
}

func (n *Node) RemoveDependency(node *Node) bool {
	if !n.dependencies[node] {
		return false
	}

	delete(n.dependencies, node)
	delete(n.dependents, n)

	return true
}
