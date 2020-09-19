package main

import (
	"reflect"
	"testing"

	"github.com/golang-collections/collections/set"
)

func TestEmptyGraph(t *testing.T) {
	g := NewGraph()
	if len(g.Nodes) != 0 {
		t.Fatalf("Expected 0 nodes, got %v\n", len(g.Nodes))
	}
}

func TestAddDirectedOneEdge(t *testing.T) {
	g := NewGraph()
	g.AddDirected("s1", "d1")

	if len(g.Nodes) != 1 {
		t.Fatalf("Expected 1 node, got %v\n", len(g.Nodes))
	}

	expected := set.New("d1")
	if !reflect.DeepEqual(g.Nodes["s1"], expected) {
		t.Fatalf("Expected %v, got %v", expected, g.Nodes["s1"])
	}
}

func TestAddDirectedTwoEdges(t *testing.T) {
	g := NewGraph()
	g.AddDirected("s1", "d1")
	g.AddDirected("s1", "d2")

	if len(g.Nodes) != 1 {
		t.Fatalf("Expected 1 node, got %v\n", len(g.Nodes))
	}

	expected := set.New("d1", "d2")
	if !reflect.DeepEqual(g.Nodes["s1"], expected) {
		t.Fatalf("Expected %v, got %v", expected, g.Nodes["s1"])
	}
}

func TestAddDirectedTwoEdgesDifferent(t *testing.T) {
	g := NewGraph()
	g.AddDirected("s1", "d1")
	g.AddDirected("s2", "d2")

	if len(g.Nodes) != 2 {
		t.Fatalf("Expected 2 nodes, got %v\n", len(g.Nodes))
	}

	expected1 := set.New("d1")
	if !reflect.DeepEqual(g.Nodes["s1"], expected1) {
		t.Fatalf("Expected %v, got %v", expected1, g.Nodes["s1"])
	}

	expected2 := set.New("d2")
	if !reflect.DeepEqual(g.Nodes["s2"], expected2) {
		t.Fatalf("Expected %v, got %v", expected2, g.Nodes["s2"])
	}
}

func TestAddUndirected(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("s1", "s2")

	if len(g.Nodes) != 2 {
		t.Fatalf("Expected 2 nodes, got %v\n", len(g.Nodes))
	}

	expected1 := set.New("s2")
	if !reflect.DeepEqual(g.Nodes["s1"], expected1) {
		t.Fatalf("Expected %v, got %v", expected1, g.Nodes["s1"])
	}

	expected2 := set.New("s1")
	if !reflect.DeepEqual(g.Nodes["s2"], expected2) {
		t.Fatalf("Expected %v, got %v", expected2, g.Nodes["s2"])
	}
}

func TestAdjacentToNotPresent(t *testing.T) {
	g := NewGraph()
	actual := g.AdjacentTo("node")

	if len(actual) != 0 {
		t.Fatalf("Expected 0 nodes, got %v", actual)
	}
}

func TestAdjacentToOneNode(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")

	actualFromA := g.AdjacentTo("a")
	expectedFromA := []string{"b"}
	if !reflect.DeepEqual(expectedFromA, actualFromA) {
		t.Fatalf("Expected %v, got %v\n", expectedFromA, actualFromA)
	}

	actualFromB := g.AdjacentTo("b")
	expectedFromB := []string{"a"}
	if !reflect.DeepEqual(expectedFromB, actualFromB) {
		t.Fatalf("Expected %v, got %v\n", expectedFromB, actualFromB)
	}
}

func TestAdjacentToTwoNodes(t *testing.T) {
	g := NewGraph()
	g.AddDirected("a", "b")
	g.AddDirected("a", "c")
	g.AddDirected("d", "f")

	actual := g.AdjacentTo("a")
	expected := []string{"b", "c"}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}
}

func TestFlattenOneVertex(t *testing.T) {
	v1 := NewVertex("vertex-1", 0)

	actual := v1.flatten()
	expected := []string{"vertex-1"}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}
}

func TestFlattenTwoVertices(t *testing.T) {
	v1 := NewVertex("vertex-1", 0)
	v2 := NewVertex("vertex-2", 1)
	v2.Parent = &v1

	actual := v2.flatten()
	expected := []string{"vertex-1", "vertex-2"}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}
}

func TestFlattenThreeVertices(t *testing.T) {
	v1 := NewVertex("vertex-1", 0)
	v2 := NewVertex("vertex-2", 1)
	v2.Parent = &v1
	v3 := NewVertex("vertex-3", 2)
	v3.Parent = &v2

	actual := v3.flatten()
	expected := []string{"vertex-1", "vertex-2", "vertex-3"}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}
}

func TestBfsTwoVertices(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")

	found, vertex := g.Bfs("a", "b", 1)
	if !found {
		t.Fatalf("Expected to find the vertex")
	}

	actual := vertex.flatten()
	expected := []string{"a", "b"}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}
}

func TestBfsThreeVertices(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")
	g.AddUndirected("b", "c")

	// a -> b
	found, vertex := g.Bfs("a", "b", 1)
	if !found {
		t.Fatalf("Expected to find the vertex")
	}

	actual := vertex.flatten()
	expected := []string{"a", "b"}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}

	// a -> b -> c
	found, vertex = g.Bfs("a", "c", 2)
	if !found {
		t.Fatalf("Expected to find the vertex")
	}

	actual = vertex.flatten()
	expected = []string{"a", "b", "c"}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}

	// a -> b -> c (but stops searching at 1)
	found, vertex = g.Bfs("a", "c", 1)
	if found {
		t.Fatalf("Expected not to find the vertex")
	}
}

func TestBfsThreeVerticesCompletelyConnected(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")
	g.AddUndirected("a", "c")
	g.AddUndirected("b", "c")

	// a -> b
	found, vertex := g.Bfs("a", "b", 1)
	if !found {
		t.Fatalf("Expected to find the vertex")
	}

	actual := vertex.flatten()
	expected := []string{"a", "b"}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}
}

func TestBfsDiamondShape(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")
	g.AddUndirected("a", "c")
	g.AddUndirected("b", "d")
	g.AddUndirected("c", "d")

	// a -> b -> d
	found, vertex := g.Bfs("a", "d", 2)
	if !found {
		t.Fatalf("Expected to find the vertex")
	}

	actual := vertex.flatten()
	expected := []string{"a", "b", "d"}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}
}

func TestBfsFourVerticesTwoConnectedComponents(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")
	g.AddUndirected("c", "d")

	found, vertex := g.Bfs("a", "b", 2)
	if !found {
		t.Fatalf("Expected to find the vertex")
	}

	actual := vertex.flatten()
	expected := []string{"a", "b"}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}

	found, vertex = g.Bfs("a", "d", 2)
	if found {
		t.Fatalf("Expected not to find the vertex")
	}
}
