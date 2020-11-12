package main

import (
	"bufio"
	"log"
	"os"
)

// ReadFileIntoSlice reads the contents of a file into a slice
func ReadFileIntoSlice(filepath string) *[]string {

	// Open the file
	file, err := os.Open(filepath)

	if err != nil {
		log.Fatalf("Failed to open %v\n", filepath)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	lines := []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return &lines
}

// FilesHaveSameContentIgnoringOrder returns true if two files have the same content, but in any order
func FilesHaveSameContentIgnoringOrder(filepath1 string, filepath2 string) bool {

	// Read the contents of the file
	contents1 := ReadFileIntoSlice(filepath1)
	contents2 := ReadFileIntoSlice(filepath2)

	// Return if the two slices are the same, ignoring the order
	return SlicesHaveSameElements(contents1, contents2)
}

// FilesHaveSameContent returns true if the files have the same content, considering the order of the rows
func FilesHaveSameContent(filepath1 string, filepath2 string) bool {

	// Read the contents of the file
	contents1 := ReadFileIntoSlice(filepath1)
	contents2 := ReadFileIntoSlice(filepath2)

	// Check the files contain the same number of rows
	if len(*contents1) != len(*contents2) {
		return false
	}

	// Walk through each row and check the contents
	for i := 0; i < len(*contents1); i++ {

		row1 := (*contents1)[i]
		row2 := (*contents2)[i]

		if row1 != row2 {
			return false
		}
	}

	// All rows must have been the same
	return true
}
