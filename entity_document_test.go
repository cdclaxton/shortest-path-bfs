package main

import (
	"reflect"
	"testing"

	"github.com/golang-collections/collections/set"
)

func TestReadEntityDocumentGraphFromFile0(t *testing.T) {
	filepath := "./test-data/entity_0.csv"
	result := ReadEntityDocumentGraphFromFile(filepath)

	expected := []EntityDocument{
		EntityDocument{
			EntityID:   "e-100",
			DocumentID: "doc-4",
		},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %v, got %v\n", expected, result)
	}
}

func TestReadEntityDocumentGraphFromFile1(t *testing.T) {
	filepath := "./test-data/entity_1.csv"
	result := ReadEntityDocumentGraphFromFile(filepath)

	expected := []EntityDocument{
		EntityDocument{
			EntityID:   "e-100",
			DocumentID: "doc-4",
		},
		EntityDocument{
			EntityID:   "e-100",
			DocumentID: "doc-1",
		},
		EntityDocument{
			EntityID:   "e-100",
			DocumentID: "doc-3",
		},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %v, got %v\n", expected, result)
	}
}

func TestReadEntityDocumentGraph(t *testing.T) {
	filepaths := []string{
		"./test-data/entity_2.csv",
		"./test-data/entity_3.csv",
	}

	result := ReadEntityDocumentGraph(filepaths)

	expected := []EntityDocument{
		EntityDocument{
			EntityID:   "e-200",
			DocumentID: "doc-1",
		},
		EntityDocument{
			EntityID:   "e-200",
			DocumentID: "doc-2",
		},
		EntityDocument{
			EntityID:   "e-300",
			DocumentID: "doc-3",
		},
		EntityDocument{
			EntityID:   "e-300",
			DocumentID: "doc-2",
		},
		EntityDocument{
			EntityID:   "e-301",
			DocumentID: "doc-4",
		},
	}

	if !reflect.DeepEqual(expected, *result) {
		t.Fatalf("Expected %v, got %v\n", expected, result)
	}
}

func TestConvertSetToSliceNoElements(t *testing.T) {
	s := set.New()

	actual := ConvertSetToSlice(s)
	expected := []string{}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}
}

func TestConvertSetToSliceOneElement(t *testing.T) {
	s := set.New()
	s.Insert("a")

	actual := ConvertSetToSlice(s)
	expected := []string{"a"}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}
}

func TestBipartiteToUnipartiteTwoEntities(t *testing.T) {
	connections := []EntityDocument{
		EntityDocument{
			EntityID:   "e-1",
			DocumentID: "d-1",
		},
		EntityDocument{
			EntityID:   "e-2",
			DocumentID: "d-1",
		},
	}

	g := BipartiteToUnipartite(&connections)

	actual1 := g.AdjacentTo("e-1")
	expected1 := []string{"e-2"}

	if !reflect.DeepEqual(expected1, actual1) {
		t.Fatalf("Expected %v, got %v\n", expected1, actual1)
	}

	actual2 := g.AdjacentTo("e-2")
	expected2 := []string{"e-1"}

	if !reflect.DeepEqual(expected2, actual2) {
		t.Fatalf("Expected %v, got %v\n", expected2, actual2)
	}
}

func TestBipartiteToUnipartiteThreeEntities(t *testing.T) {
	// Note that e-1 is connected to d-3, but d-3 isn't connected to another entity
	connections := []EntityDocument{
		EntityDocument{
			EntityID:   "e-1",
			DocumentID: "d-1",
		},
		EntityDocument{
			EntityID:   "e-2",
			DocumentID: "d-1",
		},
		EntityDocument{
			EntityID:   "e-1",
			DocumentID: "d-3",
		},
	}

	g := BipartiteToUnipartite(&connections)

	actual1 := g.AdjacentTo("e-1")
	expected1 := []string{"e-2"}

	if !reflect.DeepEqual(expected1, actual1) {
		t.Fatalf("Expected %v, got %v\n", expected1, actual1)
	}

	actual2 := g.AdjacentTo("e-2")
	expected2 := []string{"e-1"}

	if !reflect.DeepEqual(expected2, actual2) {
		t.Fatalf("Expected %v, got %v\n", expected2, actual2)
	}
}
