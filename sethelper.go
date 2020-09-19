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
