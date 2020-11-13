package main

import (
	"reflect"
	"testing"

	"github.com/golang-collections/collections/set"
)

func TestReadEntityDocumentGraphFromFile0(t *testing.T) {
	filepath := "./test/test-data/entity_0.csv"
	skipEntities := set.New()
	result := ReadEntityDocumentGraphFromFile(filepath, skipEntities)

	expected := []EntityDocument{
		EntityDocument{
			EntityID:   "e-100",
			DocumentID: "doc-4",
		},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestReadEntityDocumentGraphFromFile1(t *testing.T) {
	filepath := "./test/test-data/entity_1.csv"
	skipEntities := set.New()
	result := ReadEntityDocumentGraphFromFile(filepath, skipEntities)

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
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestTestReadEntityDocumentGraphFromFile1WithSkip(t *testing.T) {
	filepath := "./test/test-data/entity_3.csv"
	skipEntities := set.New()
	skipEntities.Insert("e-300")

	result := ReadEntityDocumentGraphFromFile(filepath, skipEntities)

	expected := []EntityDocument{
		EntityDocument{
			EntityID:   "e-301",
			DocumentID: "doc-4",
		},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestTestReadEntityDocumentGraphFromFile2WithSkip(t *testing.T) {
	filepath := "./test/test-data/entity_3.csv"
	skipEntities := set.New()
	skipEntities.Insert("e-300")
	skipEntities.Insert("e-301")

	result := ReadEntityDocumentGraphFromFile(filepath, skipEntities)

	if len(result) != 0 {
		t.Errorf("Expected list with no elements, got %v\n", len(result))
	}

}

func TestReadEntityDocumentGraph(t *testing.T) {
	filepaths := []string{
		"./test/test-data/entity_2.csv",
		"./test/test-data/entity_3.csv",
	}
	skipEntities := set.New()

	result := ReadEntityDocumentGraph(filepaths, skipEntities)

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
		t.Errorf("Expected %v, got %v\n", expected, result)
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
		t.Errorf("Expected %v, got %v\n", expected1, actual1)
	}

	actual2 := g.AdjacentTo("e-2")
	expected2 := []string{"e-1"}

	if !reflect.DeepEqual(expected2, actual2) {
		t.Errorf("Expected %v, got %v\n", expected2, actual2)
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
		t.Errorf("Expected %v, got %v\n", expected1, actual1)
	}

	actual2 := g.AdjacentTo("e-2")
	expected2 := []string{"e-1"}

	if !reflect.DeepEqual(expected2, actual2) {
		t.Errorf("Expected %v, got %v\n", expected2, actual2)
	}
}
