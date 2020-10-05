package main

import (
	"sort"

	"github.com/golang-collections/collections/set"
)

// ConvertSetToSlice converts a set to a slice
func ConvertSetToSlice(s *set.Set) []string {

	// Empty slice to hold the
	elements := []string{}

	// Walk through each element in the set
	s.Do(func(s interface{}) {
		var x string
		x = s.(string)
		elements = append(elements, x)
	})

	// Sort the elements alphabetically
	sort.Strings(elements)

	return elements
}

// SliceToSet converts a slice to a set
func SliceToSet(s []string) *set.Set {

	// Create the set
	t := set.New()

	for _, value := range s {
		t.Insert(value)
	}

	return t
}

// SetsEqual returns true if two sets have the same elements
func SetsEqual(set1 *set.Set, set2 *set.Set) bool {

	diff1 := set1.Difference(set2)
	diff2 := set2.Difference(set1)

	return diff1.Len() == 0 && diff2.Len() == 0
}

// SlicesHaveSameElements returns true if two slices have the same elements (in any order)
func SlicesHaveSameElements(s1 *[]string, s2 *[]string) bool {

	// Check the lengths
	if len(*s1) != len(*s2) {
		return false
	}

	// Convert the slices to sets
	set1 := SliceToSet(*s1)
	set2 := SliceToSet(*s2)

	// Determine if two sets are equal
	return SetsEqual(set1, set2)
}
