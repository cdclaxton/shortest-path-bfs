package main

import (
	"github.com/golang-collections/collections/queue"
	"github.com/golang-collections/collections/set"
)

// Graph represents a directed graph
type Graph struct {
	Nodes map[string]*set.Set
}

// NewGraph constructs a new Graph
func NewGraph() Graph {
	return Graph{
		Nodes: make(map[string]*set.Set),
	}
}

// AddDirected adds a directed connection in the graph
func (g *Graph) AddDirected(source string, destination string) {

	// Has the source been seen before?
	_, present := g.Nodes[source]
	if !present {
		g.Nodes[source] = set.New()
	}

	g.Nodes[source].Insert(destination)
}

// AddUndirected adds an undirected edge between source and destination vertices
func (g *Graph) AddUndirected(source string, destination string) {
	g.AddDirected(source, destination)
	g.AddDirected(destination, source)
}

// AdjacentTo returns the vertices adjacent to a given vertex
func (g *Graph) AdjacentTo(source string) []string {

	values, ok := g.Nodes[source]
	if !ok {
		return nil
	}

	return ConvertSetToSlice(values)
}

// Vertex represents a vertex in the graph
type Vertex struct {
	Identifier string
	Depth      int
	Parent     *Vertex
}

// NewVertex creates a new Vertex
func NewVertex(identifier string, depth int) Vertex {
	return Vertex{
		Identifier: identifier,
		Depth:      depth,
		Parent:     nil,
	}
}

// Flatten the vertices to a single slice
func (v *Vertex) flatten() []string {

	lineage := []string{}

	p := v
	for p != nil {
		// Prepend the lineage
		lineage = append([]string{p.Identifier}, lineage...)
		p = p.Parent
	}

	return lineage
}

// Bfs performs a Breadth First Search in the graph
func (g *Graph) Bfs(root string, goal string, maxDepth int) (bool, *Vertex) {

	// Set of the identifiers of discovered vertices
	discovered := set.New()

	// Queue to hold the vertices to visit
	q := queue.New()
	q.Enqueue(NewVertex(root, 0))

	// While there are vertices in the queue to check
	for q.Len() > 0 {

		// Take a vertex from the queue
		v := q.Dequeue().(Vertex)

		// If the vertex is the goal, then return
		if v.Identifier == goal {
			return true, &v
		}

		// Get a list of the adjacent vertices
		w := g.AdjacentTo(v.Identifier)
		newDepth := (v.Depth + 1)

		// Walk through each of the adjacent vertices
		for _, adjIdentifier := range w {

			// If the vertex hasn't been seen before
			if !discovered.Has(adjIdentifier) {

				// Add the identifier to the set of discovered identifiers
				discovered.Insert(adjIdentifier)

				if newDepth <= maxDepth {
					newVertex := NewVertex(adjIdentifier, newDepth)
					newVertex.Parent = &v
					q.Enqueue(newVertex)
				}
			}
		}
	}

	// The goal was not found
	return false, nil
}
