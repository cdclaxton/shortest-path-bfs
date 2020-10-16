package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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

// listOfKeys returns a list of keys from the Nodes data structure
func (g *Graph) listOfKeys() []string {
	keys := make([]string, len(g.Nodes))
	i := 0

	for k := range g.Nodes {
		keys[i] = k
		i++
	}

	return keys
}

// Equal determines if two graphs are equal
func (g *Graph) Equal(g2 *Graph, debug bool) bool {

	// Check the vertices
	keys1 := g.listOfKeys()
	keys2 := g2.listOfKeys()

	if !SlicesHaveSameElements(&keys1, &keys2) {
		if debug {
			fmt.Println("[!] Lists of keys are different")
			fmt.Printf("[!] Keys1: %v\n", keys1)
			fmt.Printf("[!] Keys2: %v\n", keys2)
		}
		return false
	}

	// Walk through each vertex and check its connections
	for _, vertex := range keys1 {
		conns1 := g.Nodes[vertex]
		conns2 := g2.Nodes[vertex]

		if !SetsEqual(conns1, conns2) {
			if debug {
				fmt.Printf("[!] Connections different for vertex %v", vertex)
				fmt.Printf("[!] Connections 1: %v\n", conns1)
				fmt.Printf("[!] Connections 2: %v\n", conns2)
			}
			return false
		}
	}

	return true
}

// AddDirected adds a directed connection in the graph
func (g *Graph) AddDirected(source string, destination string) {

	// Preconditions
	if source == destination {
		log.Fatalf("Source and destination vertices are identical: %v\n", source)
	}

	if len(source) == 0 {
		log.Fatal("Source vertex is empty")
	}

	if len(destination) == 0 {
		log.Fatal("Destination vertex is empty")
	}

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

	// Precondition
	if len(source) == 0 {
		log.Fatal("Source vertex is empty")
	}

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

	// Preconditions
	if len(identifier) == 0 {
		log.Fatal("Identifier is empty")
	}

	if depth < 0 {
		log.Fatalf("Invalid depth: %v\n", depth)
	}

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

// ReachableVertices finds all vertices reachable within m steps
func (g *Graph) ReachableVertices(root string, maxDepth int) (bool, *set.Set) {

	// Preconditions
	if len(root) == 0 {
		log.Fatal("Root vertex is empty")
	}

	if maxDepth < 0 {
		log.Fatalf("Maximum depth is invalid: %v\n", maxDepth)
	}

	// Set of the identifiers of discovered vertices
	discovered := set.New()
	discovered.Insert(root)

	// Check that the root vertex exists
	_, present := g.Nodes[root]
	if !present {
		return false, nil
	}

	// Queue to hold the vertices to visit
	q := queue.New()
	q.Enqueue(NewVertex(root, 0))

	// While there are vertices in the queue to check
	for q.Len() > 0 {

		// Take a vertex from the queue
		v := q.Dequeue().(Vertex)

		// Depth of any vertices adjacent to v
		newDepth := v.Depth + 1

		if newDepth <= maxDepth {

			// Get a list of the adjacent vertices
			w := g.AdjacentTo(v.Identifier)

			// Walk through each adjacent vertex
			for _, adjIdentifier := range w {

				// If the vertex hasn't been seen before
				if !discovered.Has(adjIdentifier) {

					// Add the identifier to the set of discovered identifiers
					discovered.Insert(adjIdentifier)

					newVertex := NewVertex(adjIdentifier, newDepth)
					newVertex.Parent = &v
					q.Enqueue(newVertex)

				}

			}
		}

	}

	return true, discovered
}

// Bfs performs a Breadth First Search in the graph
func (g *Graph) Bfs(root string, goal string, maxDepth int) (bool, *Vertex) {

	// Preconditions
	if len(root) == 0 {
		log.Fatal("Root vertex is empty")
	}

	if len(goal) == 0 {
		log.Fatal("Goal vertex is empty")
	}

	if maxDepth < 0 {
		log.Fatalf("Maximum depth is invalid: %v\n", maxDepth)
	}

	// Set of the identifiers of discovered vertices
	discovered := set.New()
	discovered.Insert(root)

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

		// Depth of any vertices adjacent to v
		newDepth := v.Depth + 1

		// If the adjacent vertices are within the range
		if newDepth <= maxDepth {

			// Get a list of the adjacent vertices
			w := g.AdjacentTo(v.Identifier)

			// Walk through each of the adjacent vertices
			for _, adjIdentifier := range w {

				// If the vertex hasn't been seen before
				if !discovered.Has(adjIdentifier) {

					// Add the identifier to the set of discovered identifiers
					discovered.Insert(adjIdentifier)

					// Put the vertex on the queue
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

// flattenAll flattens all of the tree nodes
func flattenAll(paths []*TreeNode) [][]string {

	flattened := [][]string{}

	for _, node := range paths {
		flattened = append(flattened, node.flatten())
	}

	return flattened
}

// AllPaths finds all the paths from root to goal up to a maximum depth
func (g *Graph) AllPaths(root string, goal string, maxDepth int) []*TreeNode {

	// Preconditions
	if len(root) == 0 {
		log.Fatal("Root vertex is empty")
	}

	if len(goal) == 0 {
		log.Fatal("Goal vertex is empty")
	}

	if maxDepth < 0 {
		log.Fatalf("Maximum depth is invalid: %v\n", maxDepth)
	}

	// Number of steps traversed from the root vertex
	numSteps := 0

	// If the goal is the root, then return without traversing the graph
	treeNode := makeTreeNode(root, root == goal)
	if treeNode.marked {
		return []*TreeNode{treeNode}
	}

	// Nodes to 'spider' from
	qCurrent := queue.New()
	qCurrent.Enqueue(treeNode)

	// Nodes to 'spider' from on the next iteration
	qNext := queue.New()

	// List of complete nodes (where goal has been found)
	complete := []*TreeNode{}

	for numSteps < maxDepth {

		for qCurrent.Len() > 0 {

			// Take a tree node from the queue representing a vertex
			node := qCurrent.Dequeue().(*TreeNode)

			if node.marked {
				log.Fatal("Trying to traverse from a marked node")
			}

			// Get a list of the adjacent vertices
			w := g.AdjacentTo(node.name)

			// Walk through each of the adjacent vertices
			for _, adjIdentifier := range w {

				if !node.containsVertex(adjIdentifier) {

					marked := adjIdentifier == goal
					child := node.makeChild(adjIdentifier, marked)

					if marked {
						complete = append(complete, child)
					} else {
						qNext.Enqueue(child)
					}
				}

			}
		}

		qCurrent = qNext
		qNext = queue.New()
		numSteps++

	}

	return complete
}

// WriteEdgeList writes the edge list to a file with the required delimiter
func (g *Graph) WriteEdgeList(filepath string, delimiter string) {

	// Precondition
	if len(delimiter) == 0 {
		log.Fatal("Delimiter is empty")
	}

	// Open the output CSV file for writing
	outputFile, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("Unable to open output file %v for writing: %v\n", filepath, err)
	}
	defer outputFile.Close()

	// Walk through the source vertices
	for source, destinations := range g.Nodes {

		// Walk through the set of destination vertices
		destinations.Do(func(s interface{}) {

			// Destination as a string
			d := s.(string)

			// Add the connection to the output file
			row := strings.Join([]string{source, d}, delimiter)
			fmt.Fprintln(outputFile, row)
		})
	}
}

// SimplifyForUndirectedGraph simplifies the graph for undirected graphs
func (g *Graph) SimplifyForUndirectedGraph() *Graph {

	// Initialise a new graph
	gUndirected := NewGraph()

	// Walk though each node
	for source, destinations := range g.Nodes {

		// Walk through each destination vertex
		destinations.Do(func(s interface{}) {

			// Destination as a string
			d := s.(string)

			// If the destination vertex comes after the source vertex then
			// add it to the simplified graph
			if d > source {
				gUndirected.AddDirected(source, d)
			}
		})
	}

	return &gUndirected
}

// WriteUndirectedEdgeList creates an edge list for an undirected graph
func (g *Graph) WriteUndirectedEdgeList(filepath string, delimiter string) {

	// Precondition
	if len(delimiter) == 0 {
		log.Fatal("Delimiter is empty")
	}

	// Simplify the graph
	simplified := g.SimplifyForUndirectedGraph()

	// Write the edge lists to a file
	simplified.WriteEdgeList(filepath, delimiter)
}
