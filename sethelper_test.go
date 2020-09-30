package main

import (
	"reflect"
	"testing"

	"github.com/golang-collections/collections/set"
)

func TestConvertSetToSliceTwoElements(t *testing.T) {
	s := set.New()
	s.Insert("a")
	s.Insert("b")

	actual := ConvertSetToSlice(s)
	expected := []string{"a", "b"}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}
}

func TestSliceToSetEmpty(t *testing.T) {
	s := []string{}
	actual := SliceToSet(s)

	expected := set.New()

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}
}

func TestSliceToSetOneElement(t *testing.T) {
	s := []string{"a"}
	actual := SliceToSet(s)

	expected := set.New()
	expected.Insert("a")

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}
}

func TestSliceToSetTwoElements(t *testing.T) {
	s := []string{"b", "c"}
	actual := SliceToSet(s)

	expected := set.New()
	expected.Insert("b")
	expected.Insert("c")

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}
}
