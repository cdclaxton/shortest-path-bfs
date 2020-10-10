package main

import (
	"log"
)

// TreeNode represents a node in a tree data structure
type TreeNode struct {
	name     string      // name of the node from the graph
	parent   *TreeNode   // parent of the node
	children []*TreeNode // children of the node
	marked   bool        // boolean flag
}

// makeTreeNode makes a new TreeNode struct
func makeTreeNode(name string, marked bool) *TreeNode {
	node := TreeNode{
		name:     name,
		parent:   nil,
		children: []*TreeNode{},
		marked:   marked,
	}

	return &node
}

// makeChild makes a child node in the tree
func (t *TreeNode) makeChild(name string, marked bool) *TreeNode {

	// Check the name is valid
	if len(name) == 0 {
		log.Fatal("Child name is empty")
	}

	// Ensure the vertex is not in the lineage
	if t.containsVertex(name) {
		log.Fatalf("Lineage already contains %v", name)
	}

	// Make the new node
	node := makeTreeNode(name, marked)
	node.parent = t

	// Add the node
	t.children = append(t.children, node)

	// Return the newly created child node
	return node
}

// containsVertex determines if a vertex or any of its parents contain a vertex
func (t *TreeNode) containsVertex(name string) bool {

	p := t
	for p != nil {
		if p.name == name {
			return true
		}
		p = p.parent
	}

	return false
}

// flatten the lineage to a slice
func (t *TreeNode) flatten() []string {

	lineage := []string{}

	p := t
	for p != nil {
		// Prepend the lineage
		lineage = append([]string{p.name}, lineage...)
		p = p.parent
	}

	return lineage
}
