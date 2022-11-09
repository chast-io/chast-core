package dependencygraph

import (
	"chast.io/core/internal/run_model/pkg/model/refactoring"
	log "github.com/sirupsen/logrus"
)

type node struct {
	self         *refactoring.Run
	dependents   map[*node]bool
	dependencies map[*node]bool
}

func newNode(run *refactoring.Run) *node {
	return &node{
		self:         run,
		dependents:   make(map[*node]bool),
		dependencies: make(map[*node]bool),
	}
}

func (n *node) addDependency(node *node) bool {
	if n.dependencies[node] {
		log.Warnf("dependency already exists: %v -> %v", n.self.ID, node.self.ID)

		return false
	}

	n.dependencies[node] = true
	node.dependents[n] = true

	return true
}

func (n *node) removeDependency(node *node) bool {
	if !n.dependencies[node] {
		return false
	}

	delete(n.dependencies, node)
	delete(node.dependents, n)

	return true
}
