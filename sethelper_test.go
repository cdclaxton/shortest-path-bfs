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
