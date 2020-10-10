package main

import (
	"reflect"
	"testing"

	"github.com/golang-collections/collections/set"
)

func TestEmptyGraph(t *testing.T) {
	g := NewGraph()
	if len(g.Nodes) != 0 {
		t.Errorf("Expected 0 nodes, got %v\n", len(g.Nodes))
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
		t.Errorf("Expected %v, got %v", expected, g.Nodes["s1"])
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
		t.Errorf("Expected %v, got %v", expected, g.Nodes["s1"])
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
		t.Errorf("Expected %v, got %v", expected1, g.Nodes["s1"])
	}

	expected2 := set.New("d2")
	if !reflect.DeepEqual(g.Nodes["s2"], expected2) {
		t.Errorf("Expected %v, got %v", expected2, g.Nodes["s2"])
	}
}

func TestAddUndirected(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("s1", "s2")

	if len(g.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %v\n", len(g.Nodes))
	}

	expected1 := set.New("s2")
	if !reflect.DeepEqual(g.Nodes["s1"], expected1) {
		t.Errorf("Expected %v, got %v", expected1, g.Nodes["s1"])
	}

	expected2 := set.New("s1")
	if !reflect.DeepEqual(g.Nodes["s2"], expected2) {
		t.Errorf("Expected %v, got %v", expected2, g.Nodes["s2"])
	}
}

func TestAdjacentToNotPresent(t *testing.T) {
	g := NewGraph()
	actual := g.AdjacentTo("node")

	if len(actual) != 0 {
		t.Errorf("Expected 0 nodes, got %v", actual)
	}
}

func TestAdjacentToOneNode(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")

	actualFromA := g.AdjacentTo("a")
	expectedFromA := []string{"b"}
	if !reflect.DeepEqual(expectedFromA, actualFromA) {
		t.Errorf("Expected %v, got %v\n", expectedFromA, actualFromA)
	}

	actualFromB := g.AdjacentTo("b")
	expectedFromB := []string{"a"}
	if !reflect.DeepEqual(expectedFromB, actualFromB) {
		t.Errorf("Expected %v, got %v\n", expectedFromB, actualFromB)
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
		t.Errorf("Expected %v, got %v\n", expected, actual)
	}
}

func TestFlattenOneVertex(t *testing.T) {
	v1 := NewVertex("vertex-1", 0)

	actual := v1.flatten()
	expected := []string{"vertex-1"}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v\n", expected, actual)
	}
}

func TestFlattenTwoVertices(t *testing.T) {
	v1 := NewVertex("vertex-1", 0)
	v2 := NewVertex("vertex-2", 1)
	v2.Parent = &v1

	actual := v2.flatten()
	expected := []string{"vertex-1", "vertex-2"}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v\n", expected, actual)
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
		t.Errorf("Expected %v, got %v\n", expected, actual)
	}
}

func TestBfsRootNodeNotPresent(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")

	found, _ := g.Bfs("c", "a", 1)
	if found {
		t.Errorf("Expected not to find the vertex")
	}
}

func TestBfsGoalNodeNotPresent(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "d")

	found, _ := g.Bfs("a", "b", 3)
	if found {
		t.Errorf("Expected not to find the vertex")
	}
}

func TestBfsTwoVertices(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")

	found, vertex := g.Bfs("a", "b", 1)
	if !found {
		t.Errorf("Expected to find the vertex")
	}

	actual := vertex.flatten()
	expected := []string{"a", "b"}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v\n", expected, actual)
	}
}

func TestBfsThreeVertices(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")
	g.AddUndirected("b", "c")

	// a -> b
	found, vertex := g.Bfs("a", "b", 1)
	if !found {
		t.Errorf("Expected to find the vertex")
	}

	actual := vertex.flatten()
	expected := []string{"a", "b"}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v\n", expected, actual)
	}

	// a -> b -> c
	found, vertex = g.Bfs("a", "c", 2)
	if !found {
		t.Errorf("Expected to find the vertex")
	}

	actual = vertex.flatten()
	expected = []string{"a", "b", "c"}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v\n", expected, actual)
	}

	// a -> b -> c (but stops searching at 1)
	found, vertex = g.Bfs("a", "c", 1)
	if found {
		t.Errorf("Expected not to find the vertex")
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
		t.Errorf("Expected to find the vertex")
	}

	actual := vertex.flatten()
	expected := []string{"a", "b"}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v\n", expected, actual)
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
		t.Errorf("Expected to find the vertex")
	}

	actual := vertex.flatten()
	expected := []string{"a", "b", "d"}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v\n", expected, actual)
	}
}

func TestBfsFourVerticesTwoConnectedComponents(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")
	g.AddUndirected("c", "d")

	found, vertex := g.Bfs("a", "b", 2)
	if !found {
		t.Errorf("Expected to find the vertex")
	}

	actual := vertex.flatten()
	expected := []string{"a", "b"}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v\n", expected, actual)
	}

	found, vertex = g.Bfs("a", "d", 2)
	if found {
		t.Errorf("Expected not to find the vertex")
	}
}

func TestReachableVerticesRootNotFound(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")

	found, _ := g.ReachableVertices("c", 1)

	if found {
		t.Errorf("Found vertex, didn't expect to\n")
	}
}

func TestReachableVerticesZeroSteps(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")
	g.AddUndirected("b", "c")
	g.AddUndirected("c", "d")

	found, actual := g.ReachableVertices("a", 0)

	if !found {
		t.Fatalf("Expected to find vertex\n")
	}

	expected := set.New()
	expected.Insert("a")

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, found %v\n", expected, actual)
	}
}

func TestReachableVerticesOneStep(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")
	g.AddUndirected("b", "c")
	g.AddUndirected("c", "d")

	found, actual := g.ReachableVertices("a", 1)

	if !found {
		t.Fatalf("Expected to find vertex\n")
	}

	expected := set.New()
	expected.Insert("a")
	expected.Insert("b")

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, found %v\n", expected, actual)
	}
}

func TestReachableVerticesTwoSteps(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")
	g.AddUndirected("b", "c")
	g.AddUndirected("c", "d")

	found, actual := g.ReachableVertices("a", 2)

	if !found {
		t.Fatalf("Expected to find vertex\n")
	}

	expected := set.New()
	expected.Insert("a")
	expected.Insert("b")
	expected.Insert("c")

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, found %v\n", expected, actual)
	}
}

func TestSimplifyForUndirectedGraph(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")

	actual := g.SimplifyForUndirectedGraph()

	expected := NewGraph()
	expected.AddDirected("a", "b")

	if !actual.Equal(&expected, false) {
		t.Errorf("Expected actual graph to be equal to the expected graph")
	}
}

func TestSimplifyForUndirectedGraph2(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")
	g.AddUndirected("a", "c")

	actual := g.SimplifyForUndirectedGraph()

	expected := NewGraph()
	expected.AddDirected("a", "b")
	expected.AddDirected("a", "c")

	if !actual.Equal(&expected, false) {
		t.Errorf("Expected actual graph to be equal to the expected graph")
	}
}

func TestSimplifyForUndirectedGraph3(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("b", "a")
	g.AddUndirected("a", "c")

	actual := g.SimplifyForUndirectedGraph()

	expected := NewGraph()
	expected.AddDirected("a", "b")
	expected.AddDirected("a", "c")

	if !actual.Equal(&expected, false) {
		t.Errorf("Expected actual graph to be equal to the expected graph")
	}
}

func TestSimplifyForUndirectedGraph4(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")
	g.AddUndirected("a", "c")
	g.AddUndirected("d", "a")

	actual := g.SimplifyForUndirectedGraph()

	expected := NewGraph()
	expected.AddDirected("a", "b")
	expected.AddDirected("a", "c")
	expected.AddDirected("a", "d")

	if !actual.Equal(&expected, false) {
		t.Errorf("Expected actual graph to be equal to the expected graph")
	}
}

func TestWriteEdgeList(t *testing.T) {

	// Construct a test directed graph
	g := NewGraph()
	g.AddDirected("a", "b")
	g.AddDirected("a", "c")
	g.AddDirected("b", "d")
	g.AddDirected("d", "e")

	// Write the graph to a text file
	actualFilepath := "./test-writing/actual-output-1.csv"
	g.WriteEdgeList(actualFilepath, ",")

	// Check the result
	if !FilesHaveSameContentIgnoringOrder(actualFilepath, "./test-writing/expected-output-1.csv") {
		t.Fatalf("Actual results differ from expected results\n")
	}
}

func TestWriteUndirectedEdgeList(t *testing.T) {

	// Construct a test undirected graph
	g := NewGraph()
	g.AddUndirected("b", "a")
	g.AddUndirected("b", "d")
	g.AddUndirected("c", "b")
	g.AddUndirected("c", "e")
	g.AddUndirected("c", "f")

	// Write the graph to a text file
	actualFilepath := "./test-writing/actual-output-2.csv"
	g.WriteUndirectedEdgeList(actualFilepath, ",")

	// Check the result
	if !FilesHaveSameContentIgnoringOrder(actualFilepath, "./test-writing/expected-output-2.csv") {
		t.Fatalf("Actual results differ from expected results\n")
	}
}

func TestAllPaths2Vertices(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")

	paths := g.AllPaths("a", "b", 1)
	actualPaths := flattenAll(paths)

	expectedPaths := [][]string{
		{"a", "b"},
	}

	if !reflect.DeepEqual(expectedPaths, actualPaths) {
		t.Errorf("Expected %v, got %v", expectedPaths, actualPaths)
	}
}

func TestAllPaths3Vertices(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")
	g.AddUndirected("b", "c")

	// Stop too early
	pathsStopped := g.AllPaths("a", "c", 1)
	if len(pathsStopped) > 0 {
		t.Errorf("Didn't expect a path, found %v paths", len(pathsStopped))
	}

	// Stop after 2 steps
	paths := g.AllPaths("a", "c", 2)
	actualPaths := flattenAll(paths)

	expectedPaths := [][]string{
		{"a", "b", "c"},
	}

	if !reflect.DeepEqual(expectedPaths, actualPaths) {
		t.Errorf("Expected %v, got %v", expectedPaths, actualPaths)
	}
}

func TestAllPaths4Vertices(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")
	g.AddUndirected("b", "d")
	g.AddUndirected("a", "c")
	g.AddUndirected("c", "d")

	// Stop too early
	pathsStopped := g.AllPaths("a", "d", 1)
	if len(pathsStopped) > 0 {
		t.Errorf("Didn't expect a path, found %v paths", len(pathsStopped))
	}

	// Stop after 2 steps
	paths := g.AllPaths("a", "d", 2)
	actualPaths := flattenAll(paths)

	expectedPaths := [][]string{
		{"a", "b", "d"},
		{"a", "c", "d"},
	}

	if !reflect.DeepEqual(expectedPaths, actualPaths) {
		t.Errorf("Expected %v, got %v", expectedPaths, actualPaths)
	}
}

func TestAllPaths6Vertices(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")
	g.AddUndirected("b", "c")
	g.AddUndirected("b", "d")
	g.AddUndirected("c", "e")
	g.AddUndirected("d", "e")
	g.AddUndirected("e", "f")

	paths := g.AllPaths("a", "e", 4)
	actualPaths := flattenAll(paths)

	expectedPaths := [][]string{
		{"a", "b", "c", "e"},
		{"a", "b", "d", "e"},
	}

	if !reflect.DeepEqual(expectedPaths, actualPaths) {
		t.Errorf("Expected %v, got %v", expectedPaths, actualPaths)
	}
}

func TestAllPaths6Vertices2(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")
	g.AddUndirected("b", "c")
	g.AddUndirected("b", "d")
	g.AddUndirected("c", "d")
	g.AddUndirected("c", "e")
	g.AddUndirected("d", "e")
	g.AddUndirected("e", "f")
	g.AddUndirected("d", "f")

	paths := g.AllPaths("a", "f", 4)
	actualPaths := flattenAll(paths)

	expectedPaths := [][]string{
		{"a", "b", "d", "f"},
		{"a", "b", "c", "d", "f"},
		{"a", "b", "c", "e", "f"},
		{"a", "b", "d", "e", "f"},
	}

	if !reflect.DeepEqual(expectedPaths, actualPaths) {
		t.Errorf("Expected %v, got %v", expectedPaths, actualPaths)
	}
}

func TestAllPaths8Vertices(t *testing.T) {
	g := NewGraph()
	g.AddUndirected("a", "b")
	g.AddUndirected("b", "c")
	g.AddUndirected("c", "e")
	g.AddUndirected("c", "f")
	g.AddUndirected("e", "h")
	g.AddUndirected("f", "h")
	g.AddUndirected("b", "d")
	g.AddUndirected("d", "g")
	g.AddUndirected("g", "h")

	paths := g.AllPaths("a", "h", 4)
	actualPaths := flattenAll(paths)

	expectedPaths := [][]string{
		{"a", "b", "c", "e", "h"},
		{"a", "b", "c", "f", "h"},
		{"a", "b", "d", "g", "h"},
	}

	if !reflect.DeepEqual(expectedPaths, actualPaths) {
		t.Errorf("Expected %v, got %v", expectedPaths, actualPaths)
	}
}
