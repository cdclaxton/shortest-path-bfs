package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/golang-collections/collections/set"
)

// EntityDocument represents an entity-document relationship
type EntityDocument struct {
	EntityID   string
	DocumentID string
}

// ReadEntityDocumentGraphFromFile reads entity-document relationships from a file
func ReadEntityDocumentGraphFromFile(filepath string, skipEntities *set.Set) []EntityDocument {

	fmt.Printf("[>] Reading entity-document data from: %v\n", filepath)

	// Open the file for reading
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Couldn't open CSV file ", err)
	}

	// Ensure the file is closed
	defer file.Close()

	// Initialise the slice of connections
	var connections []EntityDocument

	// Parse the file
	r := csv.NewReader(file)
	numRowsRead := 0

	for {

		// Read a row from the file
		row, err := r.Read()
		numRowsRead++

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("Error reading CSV file: ", err)
		}

		// Ignore the header
		if numRowsRead == 1 {
			continue
		}

		if len(row) != 2 {
			log.Fatal("Invalid row: ", row)
		}

		docEnt := EntityDocument{
			EntityID:   row[0],
			DocumentID: row[1],
		}

		if !skipEntities.Has(docEnt.EntityID) {
			connections = append(connections, docEnt)
		}

	}

	fmt.Printf("[>] Read %v rows from file %v\n", numRowsRead, filepath)

	return connections
}

// ReadEntityDocumentGraph reads the entity-document graph from file
func ReadEntityDocumentGraph(files []string, skipEntities *set.Set) *[]EntityDocument {

	var allConnections []EntityDocument

	// Read the connections from each file
	for _, filePath := range files {
		conns := ReadEntityDocumentGraphFromFile(filePath, skipEntities)
		allConnections = append(allConnections, conns...)
	}

	return &allConnections
}

// BipartiteToUnipartite converts a bipartite graph to a unipartite graph
func BipartiteToUnipartite(connections *[]EntityDocument) *Graph {

	// Map of document IDs to a set of entity IDs
	docToEntities := make(map[string]*set.Set)

	for _, conn := range *connections {

		// If the document ID hasn't been seen, then add an empty Set
		_, present := docToEntities[conn.DocumentID]
		if !present {
			docToEntities[conn.DocumentID] = set.New()
		}

		docToEntities[conn.DocumentID].Insert(conn.EntityID)
	}

	// Build the graph of just entities
	g := NewGraph()

	for _, entIDs := range docToEntities {
		if entIDs.Len() == 2 {
			elements := ConvertSetToSlice(entIDs)
			g.AddUndirected(elements[0], elements[1])
		} else if entIDs.Len() > 2 {
			fmt.Printf("[!] Expected at most 2 links between document, found %v\n", entIDs.Len())
		}
	}

	return &g
}
