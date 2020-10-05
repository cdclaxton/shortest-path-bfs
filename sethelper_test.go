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
		t.Errorf("Expected %v, got %v\n", expected, actual)
	}
}

func TestSliceToSetEmpty(t *testing.T) {
	s := []string{}
	actual := SliceToSet(s)

	expected := set.New()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v\n", expected, actual)
	}
}

func TestSliceToSetOneElement(t *testing.T) {
	s := []string{"a"}
	actual := SliceToSet(s)

	expected := set.New()
	expected.Insert("a")

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v\n", expected, actual)
	}
}

func TestSliceToSetTwoElements(t *testing.T) {
	s := []string{"b", "c"}
	actual := SliceToSet(s)

	expected := set.New()
	expected.Insert("b")
	expected.Insert("c")

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v\n", expected, actual)
	}
}

type SlicesEqualTestCase struct {
	slice1   []string
	slice2   []string
	expected bool
}

var slicesEqualTestCases = []SlicesEqualTestCase{
	{
		slice1:   []string{"a"},
		slice2:   []string{"a"},
		expected: true,
	},
	{
		slice1:   []string{"a"},
		slice2:   []string{"b"},
		expected: false,
	},
	{
		slice1:   []string{"a", "b"},
		slice2:   []string{"a", "b"},
		expected: true,
	},
	{
		slice1:   []string{"a", "b"},
		slice2:   []string{"b", "a"},
		expected: true,
	},
	{
		slice1:   []string{"a", "b"},
		slice2:   []string{"a", "c"},
		expected: false,
	},
}

func TestSlicesHaveSameElements(t *testing.T) {

	for _, test := range slicesEqualTestCases {
		actual := SlicesHaveSameElements(&test.slice1, &test.slice2)
		if actual != test.expected {
			t.Errorf("Expected %v, got %v for test case %v", test.expected, actual, test)
		}
	}
}
