package main

import (
	"reflect"
	"testing"
)

func TestMakeTreeNode(t *testing.T) {
	actual := makeTreeNode("node-1", false)

	expected := TreeNode{
		name:     "node-1",
		parent:   nil,
		children: []*TreeNode{},
		marked:   false,
	}

	if !reflect.DeepEqual(expected, *actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestMakeChild1(t *testing.T) {
	a := makeTreeNode("a", false)
	b := a.makeChild("b", true)

	if b.parent != a {
		t.Errorf("Expected 'a' to be parent, got %v", b.parent)
	}

	if !b.marked {
		t.Error("Expected 'b' to be marked")
	}

	if len(a.children) != 1 {
		t.Errorf("Expected 'a' to have one child, got %v", len(a.children))
	}

	expectedChildrenA := []*TreeNode{b}
	if !reflect.DeepEqual(expectedChildrenA, a.children) {
		t.Errorf("Expected %v, got %v", expectedChildrenA, a.children)
	}
}

func TestMakeChild2(t *testing.T) {
	a := makeTreeNode("a", false)
	b := a.makeChild("b", true)
	c := a.makeChild("c", false)

	if b.parent != a {
		t.Errorf("Expected 'a' to be parent, got %v", b.parent)
	}

	if !b.marked {
		t.Error("Expected 'b' to be marked")
	}

	if len(b.children) > 0 {
		t.Errorf("Expected 'b' to have 0 children, got %v", len(b.children))
	}

	if c.parent != a {
		t.Errorf("Expected 'a' to be parent, got %v", b.parent)
	}

	if c.marked {
		t.Error("Expected 'c' not to be marked")
	}

	if len(c.children) > 0 {
		t.Errorf("Expected 'c' to have 0 children, got %v", len(c.children))
	}

	expectedChildrenA := []*TreeNode{b, c}
	if !reflect.DeepEqual(expectedChildrenA, a.children) {
		t.Errorf("Expected %v, got %v", expectedChildrenA, a.children)
	}
}

func TestContainsVertex1(t *testing.T) {
	a := makeTreeNode("a", false)

	if !a.containsVertex("a") {
		t.Error("Expected 'a' to contain 'a'")
	}
}

func TestContainsVertex2(t *testing.T) {
	a := makeTreeNode("a", false)
	b := a.makeChild("b", true)

	if !a.containsVertex("a") {
		t.Error("Expected 'a' to contain 'a'")
	}

	if a.containsVertex("b") {
		t.Error("Expected 'a' not to contain 'b'")
	}

	if !b.containsVertex("a") {
		t.Error("Expected 'b' to contain 'a'")
	}

	if !b.containsVertex("b") {
		t.Error("Expected 'b' to contain 'b'")
	}
}

func TestFlatten1(t *testing.T) {
	a := makeTreeNode("a", false)

	actual := a.flatten()
	expected := []string{"a"}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestFlatten2(t *testing.T) {
	a := makeTreeNode("a", false)
	b := a.makeChild("b", false)
	a.makeChild("c", false)

	actualA := a.flatten()
	expectedA := []string{"a"}

	if !reflect.DeepEqual(expectedA, actualA) {
		t.Errorf("Expected %v, got %v", expectedA, actualA)
	}

	actualB := b.flatten()
	expectedB := []string{"a", "b"}

	if !reflect.DeepEqual(expectedB, actualB) {
		t.Errorf("Expected %v, got %v", expectedB, actualB)
	}
}

func TestComplexCase(t *testing.T) {

	// Tier 0
	a := makeTreeNode("a", false)

	// Tier 1
	b := a.makeChild("b", false)

	// Tier 2
	c := b.makeChild("c", false)
	d := b.makeChild("d", false)

	// Tier 3
	c.makeChild("e", false)
	f := c.makeChild("f", false)
	d.makeChild("g", false)
	h := d.makeChild("h", false)
	j1 := d.makeChild("j", true)

	// Tier 4
	j2 := f.makeChild("j", true)
	j3 := h.makeChild("j", true)

	actualPathJ1 := j1.flatten()
	expectedPathJ1 := []string{"a", "b", "d", "j"}

	if !reflect.DeepEqual(expectedPathJ1, actualPathJ1) {
		t.Errorf("Expected %v, got %v", expectedPathJ1, actualPathJ1)
	}

	actualPathJ2 := j2.flatten()
	expectedPathJ2 := []string{"a", "b", "c", "f", "j"}

	if !reflect.DeepEqual(expectedPathJ2, actualPathJ2) {
		t.Errorf("Expected %v, got %v", expectedPathJ2, actualPathJ2)
	}

	actualPathJ3 := j3.flatten()
	expectedPathJ3 := []string{"a", "b", "d", "h", "j"}

	if !reflect.DeepEqual(expectedPathJ3, actualPathJ3) {
		t.Errorf("Expected %v, got %v", expectedPathJ3, actualPathJ3)
	}
}
