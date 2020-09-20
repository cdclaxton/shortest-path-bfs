package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// EntityPairsConfig represents the entity pairs to find paths for
type EntityPairsConfig struct {
	PairsToFind   []string `json:"pairs"`
	PairDelimiter string   `json:"delimiter"`
}

// OutputConfig represents the config for the output from the BFS
type OutputConfig struct {
	MaxDepth        int    `json:"max_depth"`
	OutputFile      string `json:"output_file"`
	OutputDelimiter string `json:"delimiter"`
	PathDelimiter   string `json:"path_delimiter"`
	WebAppLink      string `json:"webapp_link"`
}

// PathConfig represents the JSON config
type PathConfig struct {
	InputFiles  []string          `json:"input_files"`
	EntityPairs EntityPairsConfig `json:"entity_pairs"`
	Output      OutputConfig      `json:"output"`
}

// display the path config
func (c *PathConfig) display() {
	fmt.Println("    Number of input files:   ", len(c.InputFiles))
	fmt.Println("    Number of paths to find: ", len(c.EntityPairs.PairsToFind))
	fmt.Println("    Pair delimiter:          ", c.EntityPairs.PairDelimiter)
	fmt.Println("    Maximum depth:           ", c.Output.MaxDepth)
	fmt.Println("    Output file:             ", c.Output.OutputFile)
	fmt.Println("    Delimiter:               ", c.Output.OutputDelimiter)
	fmt.Println("    Path delimiter:          ", c.Output.PathDelimiter)
	fmt.Println("    Web-app link template:   ", c.Output.WebAppLink)
}

// readConfig reads the JSON configuration from a file
func readConfig(filePath string) PathConfig {

	// Read the contents of the file
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Unable to read config file: %v", filePath)
	}

	// Unmarshall the JSON in the config file
	config := PathConfig{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatal("Unable to unmarshall JSON")
	}

	return config
}

// PathResult represents a shortest path
type PathResult struct {
	SourceEntityID      string
	DestinationEntityID string
	NumberOfHops        int
	Path                []string
	WebAppLink          string
}

// buildWebAppLink builds the web-app link
func buildWebAppLink(template string, path []string) string {

	// List of comma-separated entity IDs
	entityIds := strings.Join(path, ",")

	// Use the template to build the URL
	return strings.Replace(template, "<ENTITY_IDS>", entityIds, -1)
}

// NewPathResult returns a PathResult based on a list of vertices
func NewPathResult(source string, destination string, vertices []string, webAppTemplate string) PathResult {
	return PathResult{
		SourceEntityID:      source,
		DestinationEntityID: destination,
		NumberOfHops:        len(vertices) - 1,
		Path:                vertices,
		WebAppLink:          buildWebAppLink(webAppTemplate, vertices),
	}
}

func (r *PathResult) display() string {
	return fmt.Sprintf("%v -> %v (%v hops): %v",
		r.SourceEntityID, r.DestinationEntityID, r.NumberOfHops, r.Path)
}

func (r *PathResult) toString(delimiter string, pathDelimiter string) string {

	// Build a representation of the path as a string
	path := strings.Join(r.Path, pathDelimiter)

	parts := []string{
		r.SourceEntityID,
		r.DestinationEntityID,
		strconv.Itoa(r.NumberOfHops),
		path,
		r.WebAppLink,
	}

	// Join the elements and return
	return strings.Join(parts, delimiter)
}

func pathResultHeader(delimiter string) string {
	parts := []string{
		"Source entity ID",
		"Destination entity ID",
		"Number of hops",
		"Path",
		"Link",
	}

	// Join the elements and return
	return strings.Join(parts, delimiter)
}

// extractEntityPair parses the entity pair
func extractEntityPair(pair string, delimiter string) (string, string, error) {

	// Split the pair, e.g. "e-1,e-2" into entities
	parts := strings.Split(pair, delimiter)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("Expected 2 entity IDs, got %v in %v", len(parts), pair)
	}

	return parts[0], parts[1], nil
}

// Perform BFS
func performBfs(g *Graph, entityConfig EntityPairsConfig, outputConfig OutputConfig) {

	// Open the output CSV file for writing
	outputFile, err := os.Create(outputConfig.OutputFile)
	if err != nil {
		log.Fatalf("Unable to open output file %v for writing: %v", outputConfig.OutputFile, err)
	}
	defer outputFile.Close()

	// Write the header
	fmt.Fprintln(outputFile, pathResultHeader(outputConfig.OutputDelimiter))

	// Walk through each pair of entity IDs
	for pairIndex, pair := range entityConfig.PairsToFind {

		// Provide feedback on long-running jobs
		if pairIndex%10000 == 0 {
			fmt.Printf("[>] Processed %v pairs of %v\n", pairIndex+1, len(entityConfig.PairsToFind))
		}

		// Extract the two entity IDs
		e1, e2, err := extractEntityPair(pair, entityConfig.PairDelimiter)
		if err != nil {
			fmt.Printf("[!] Failed to parse line. %v", err)
		}

		//fmt.Printf("[>] Checking %v -> %v\n", e1, e2)

		// Compute the shortest path using BFS
		found, vertex := g.Bfs(e1, e2, outputConfig.MaxDepth)

		// If the path could be found, add it to the output file
		if found {

			// Build the PathResult
			result := NewPathResult(e1, e2, vertex.flatten(), outputConfig.WebAppLink)

			// Display the result
			fmt.Printf("[>] %v\n", result.display())

			// Add the result to the file
			fmt.Fprintln(outputFile, result.toString(outputConfig.OutputDelimiter, outputConfig.PathDelimiter))
		}

	}

}

// PerformBfsFromConfig performs BFS based on a config file
func PerformBfsFromConfig(configFilepath string) {

	// Read the JSON configuration
	fmt.Println("[>] Reading configuration ...")
	config := readConfig(configFilepath)
	config.display()

	// Read the entity-document relationships from file
	fmt.Println("[>] Reading entity-document graph from file ...")
	connections := ReadEntityDocumentGraph(config.InputFiles)

	// Convert the bipartite graph to a unipartite graph
	graph := BipartiteToUnipartite(connections)
	fmt.Printf("[>] Graph has %v vertices\n", len(graph.Nodes))

	// Perform BFS
	fmt.Printf("[>] Performing BFS on %v vertex pairs\n", len(config.EntityPairs.PairsToFind))
	performBfs(graph, config.EntityPairs, config.Output)
}

func main() {
	println("Shortest path calculator using a bipartite to unipartite transformation and the Breadth First Search algorithm")
	PerformBfsFromConfig("./config.json")
}
