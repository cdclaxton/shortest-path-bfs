package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/golang-collections/collections/set"
)

// EntityDocument represents an entity-document relationship
type EntityDocument struct {
	EntityID   string // entity ID
	DocumentID string // document ID
}

// ReadEntityDocumentGraphFromFile reads entity-document relationships from a file, skipping the required entities
func ReadEntityDocumentGraphFromFile(filepath string, skipEntities *set.Set) []EntityDocument {

	log.Printf("Reading entity-document data from: %v\n", filepath)

	// Open the file for reading
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("[!] Couldn't open CSV file ", err)
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
			log.Fatal("[!] Error reading CSV file: ", err)
		}

		// Ignore the header
		if numRowsRead == 1 {
			continue
		}

		if len(row) != 2 {
			log.Fatal("[!] Invalid row: ", row)
		}

		docEnt := EntityDocument{
			EntityID:   row[0],
			DocumentID: row[1],
		}

		if !skipEntities.Has(docEnt.EntityID) {
			connections = append(connections, docEnt)
		}

	}

	log.Printf("Read %v rows from file %v\n", numRowsRead, filepath)

	return connections
}

// ReadEntityDocumentGraph reads the entity-document graph from a list of files
func ReadEntityDocumentGraph(files []string, skipEntities *set.Set) *[]EntityDocument {

	var allConnections []EntityDocument

	// Read the connections from each file
	for _, filePath := range files {
		conns := ReadEntityDocumentGraphFromFile(filePath, skipEntities)
		allConnections = append(allConnections, conns...)
	}

	return &allConnections
}

// BipartiteToUnipartite converts a bipartite graph to a unipartite graph by collapsing document links
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

	// Number of documents connecting the set number of entities
	numOneEntity := 0
	numTwoEntities := 0
	numThreeEntities := 0
	numFourOrMoreEntities := 0

	for _, entIDs := range docToEntities {
		if entIDs.Len() == 1 {
			numOneEntity++
		} else if entIDs.Len() == 2 {
			elements := ConvertSetToSlice(entIDs)
			g.AddUndirected(elements[0], elements[1])
			numTwoEntities++
		} else if entIDs.Len() == 3 {
			// This case isn't expected, but is handled
			elements := ConvertSetToSlice(entIDs)
			g.AddUndirected(elements[0], elements[1])
			g.AddUndirected(elements[0], elements[2])
			g.AddUndirected(elements[1], elements[2])
			numThreeEntities++
		} else {
			numFourOrMoreEntities++
		}
	}

	log.Printf("Summary - Number of documents with 1 entity:   %v\n", numOneEntity)
	log.Printf("Summary - Number of documents with 2 entities: %v\n", numTwoEntities)
	log.Printf("Summary - Number of documents with 3 entities: %v\n", numThreeEntities)
	log.Printf("Summary - Number of documents with 4+ entity:  %v\n", numFourOrMoreEntities)

	return &g
}
