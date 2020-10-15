package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// DataSource represents a named data source with entity IDs
type DataSource struct {
	Name      string   `json:"name"`
	EntityIds []string `json:"entity_ids"`
}

// EntityConfig represents the entity pairs for which to find paths
type EntityConfig struct {
	DataSources []DataSource `json:"data_sources"`
	Skip        []string     `json:"skip"`
}

// OutputConfig represents the config for the output from the BFS
type OutputConfig struct {
	MaxDepth        int    `json:"max_depth"`
	FindAllPaths    bool   `json:"find_all_paths"`
	OutputFile      string `json:"output_file"`
	OutputDelimiter string `json:"delimiter"`
	PathDelimiter   string `json:"path_delimiter"`
	WebAppLink      string `json:"webapp_link"`
	UnipartiteFile  string `json:"unipartite"`
}

// PathConfig represents the JSON config
type PathConfig struct {
	InputFiles []string     `json:"input_files"`
	Entities   EntityConfig `json:"entities"`
	Output     OutputConfig `json:"output"`
}

// display the path config
func (c *PathConfig) display() {
	fmt.Println("    Number of input files:      ", len(c.InputFiles))
	fmt.Println("    Number of data sources:     ", len(c.Entities.DataSources))
	fmt.Println("    Number of entities to skip: ", len(c.Entities.Skip))
	fmt.Println("    Maximum depth:              ", c.Output.MaxDepth)
	fmt.Println("    Find all paths:             ", c.Output.FindAllPaths)
	fmt.Println("    Output file:                ", c.Output.OutputFile)
	fmt.Println("    Delimiter:                  ", c.Output.OutputDelimiter)
	fmt.Println("    Path delimiter:             ", c.Output.PathDelimiter)
	fmt.Println("    Web-app link template:      ", c.Output.WebAppLink)
	fmt.Println("    Unipartite graph file:      ", c.Output.UnipartiteFile)
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
	SourceEntityID              string
	SourceEntityDataSource      string
	DestinationEntityID         string
	DestinationEntityDataSource string
	NumberOfHops                int
	Path                        []string
	WebAppLink                  string
}

// buildWebAppLink builds the web-app link
func buildWebAppLink(template string, path []string) string {

	// List of comma-separated entity IDs
	entityIds := strings.Join(path, ",")

	// Use the template to build the URL
	return strings.Replace(template, "<ENTITY_IDS>", entityIds, -1)
}

// NewPathResult returns a PathResult based on a list of vertices
func NewPathResult(source string, sourceDataSource string,
	destination string, destinationDataSource string,
	vertices []string, webAppTemplate string) PathResult {

	// Check the parameters
	if len(source) == 0 {
		log.Fatal("Source entity ID is blank")
	}

	if len(sourceDataSource) == 0 {
		log.Fatal("Data source for the source entity is blank")
	}

	if len(destination) == 0 {
		log.Fatal("Destination entity ID is blank")
	}

	if len(destinationDataSource) == 0 {
		log.Fatal("Data source for the destination entity is blank")
	}

	if len(vertices) < 2 {
		log.Fatalf("List of vertices on path is too small (%v)", len(vertices))
	}

	return PathResult{
		SourceEntityID:              source,
		SourceEntityDataSource:      sourceDataSource,
		DestinationEntityID:         destination,
		DestinationEntityDataSource: destinationDataSource,
		NumberOfHops:                len(vertices) - 1,
		Path:                        vertices,
		WebAppLink:                  buildWebAppLink(webAppTemplate, vertices),
	}
}

// display produces a string representation of the path for stdout
func (r *PathResult) display() string {
	return fmt.Sprintf("%v:%v -> %v:%v (%v hops) %v",
		r.SourceEntityID, r.SourceEntityDataSource,
		r.DestinationEntityID, r.DestinationEntityDataSource,
		r.NumberOfHops, r.Path)
}

// toString converts a path result to delimited form for writing to file
func (r *PathResult) toString(delimiter string, pathDelimiter string) string {

	// Build a representation of the path as a simple delimited string
	path := strings.Join(r.Path, pathDelimiter)

	parts := []string{
		r.SourceEntityID,
		r.SourceEntityDataSource,
		r.DestinationEntityID,
		r.DestinationEntityDataSource,
		strconv.Itoa(r.NumberOfHops),
		path,
		r.WebAppLink,
	}

	// Join the elements and return
	return strings.Join(parts, delimiter)
}

// pathResultHeader returns the header for the delimited file
func pathResultHeader(delimiter string) string {
	parts := []string{
		"Source entity ID",
		"Source entity data source",
		"Destination entity ID",
		"Destination entity data source",
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
		return "", "", fmt.Errorf("[!] Expected 2 entity IDs, got %v in %v", len(parts), pair)
	}

	return parts[0], parts[1], nil
}

// findAndRecordShortestPaths finds the shortest path and writes to file
func findAndRecordShortestPaths(g *Graph,
	source string, sourceDataSource string,
	destination string, destinationDataSource string,
	outputConfig OutputConfig, outputFile *os.File) {

	if outputConfig.FindAllPaths {

		// Find all the paths between the source and destination up to a maximum length
		paths := g.AllPaths(source, destination, outputConfig.MaxDepth)

		if len(paths) == 0 {
			fmt.Printf("[!] Vertex %v was deemed reachable from %v, but no path!\n", destination, source)
		} else {
			for _, path := range paths {
				result := NewPathResult(source, sourceDataSource,
					destination, destinationDataSource,
					path.flatten(), outputConfig.WebAppLink)
				fmt.Printf("[>] %v\n", result.display())
				fmt.Fprintln(outputFile, result.toString(outputConfig.OutputDelimiter, outputConfig.PathDelimiter))
			}
		}

	} else {
		// Compute the shortest path using BFS
		found, vertex := g.Bfs(source, destination, outputConfig.MaxDepth)

		if !found {
			fmt.Printf("[!] Vertex %v was deemed reachable from %v, but no path!\n", destination, source)
		} else {

			// Build the PathResult
			result := NewPathResult(source, sourceDataSource,
				destination, destinationDataSource,
				vertex.flatten(), outputConfig.WebAppLink)

			// Display the result
			fmt.Printf("[>] %v\n", result.display())

			// Add the result to the file
			fmt.Fprintln(outputFile, result.toString(outputConfig.OutputDelimiter, outputConfig.PathDelimiter))
		}
	}

}

// totalNumberOfPairs returns the total number of pairs of entities
func totalNumberOfPairs(dataSources *[]DataSource) int {

	if len(*dataSources) < 2 {
		return 0
	}

	total := 0

	// Walk through each ordered pair of data sources
	for i := 0; i < len(*dataSources)-1; i++ {
		for j := i + 1; j < len(*dataSources); j++ {

			lenA := len((*dataSources)[i].EntityIds)
			lenB := len((*dataSources)[j].EntityIds)

			total += (lenA * lenB)

		}
	}

	return total
}

// performBfs performs breadth first search or exhaustive search given a graph and config
func performBfs(g *Graph, entityConfig EntityConfig, outputConfig OutputConfig) {

	// Open the output CSV file for writing
	outputFile, err := os.Create(outputConfig.OutputFile)
	if err != nil {
		log.Fatalf("[!] Unable to open output file %v for writing: %v\n", outputConfig.OutputFile, err)
	}
	defer outputFile.Close()

	// Write the header
	fmt.Fprintln(outputFile, pathResultHeader(outputConfig.OutputDelimiter))

	// Total number of entity pairs to check
	totalPairs := totalNumberOfPairs(&entityConfig.DataSources)
	numPairsProcessed := 0

	// Make a set of entities to skip
	skipEntities := SliceToSet(entityConfig.Skip)
	numEntitiesSkipped := 0

	// Walk through all pairs of data sources
	for i := 0; i < len(entityConfig.DataSources)-1; i++ {
		for j := i + 1; j < len(entityConfig.DataSources); j++ {

			fmt.Printf("[>] Checking connections for data sources %v <--> %v\n",
				entityConfig.DataSources[i].Name,
				entityConfig.DataSources[j].Name)

			for _, source := range entityConfig.DataSources[i].EntityIds {

				// Skip the entity if required
				if skipEntities.Has(source) {
					numPairsProcessed += len(entityConfig.DataSources[i].EntityIds)
					numEntitiesSkipped++
					continue
				}

				// Set of all vertices within reach of the source vertex
				found, reachable := g.ReachableVertices(source, outputConfig.MaxDepth)

				// If the source vertex was not found, just continue to the next vertex
				if !found {
					numPairsProcessed += len(entityConfig.DataSources[i].EntityIds)
					continue
				}

				for _, destination := range entityConfig.DataSources[j].EntityIds {

					// Skip the entity if required
					if skipEntities.Has(destination) {
						numPairsProcessed++
						numEntitiesSkipped++
						continue
					}

					// Provide feedback on long-running jobs
					if numPairsProcessed%10000 == 0 {
						fmt.Printf("[>] Processed %v pairs of %v\n", numPairsProcessed, totalPairs)
					}

					// If the destination is reachable from the source, then find and record the shortest path
					if reachable.Has(destination) {
						findAndRecordShortestPaths(g,
							source,
							entityConfig.DataSources[i].Name,
							destination,
							entityConfig.DataSources[j].Name,
							outputConfig,
							outputFile)
					}

					numPairsProcessed++
				}

			}

		}

		fmt.Printf("[>] Number of entities skipped: %v\n", numEntitiesSkipped)
	}

}

// PerformBfsFromConfig performs BFS based on a config file
func PerformBfsFromConfig(configFilepath string) {

	// Read the JSON configuration
	t0 := time.Now()
	fmt.Println("[>] Reading configuration ...")
	config := readConfig(configFilepath)
	config.display()

	// Check there are at least two data sources to find connections
	if len(config.Entities.DataSources) < 2 {
		fmt.Println("[!] At least two data sources must be specified in the config")
		return
	}

	// Read the entity-document relationships from file
	fmt.Println("[>] Reading entity-document graph from file ...")
	t1 := time.Now()
	connections := ReadEntityDocumentGraph(config.InputFiles, SliceToSet(config.Entities.Skip))
	fmt.Printf("[>] Entity-document graph read in %v\n", time.Now().Sub(t1))

	// Convert the bipartite graph to a unipartite graph
	t2 := time.Now()
	graph := BipartiteToUnipartite(connections)
	fmt.Printf("[>] Graph has %v vertices\n", len(graph.Nodes))
	fmt.Printf("[>] Bipartite to unipartite conversion completed in %v\n", time.Now().Sub(t2))

	// Write the unipartite graph to file (if required)
	if len(config.Output.UnipartiteFile) > 0 {
		fmt.Printf("[>] Writing unipartite graph to file: %v\n", config.Output.UnipartiteFile)
		graph.WriteUndirectedEdgeList(config.Output.UnipartiteFile, config.Output.PathDelimiter)
	}

	// Perform shortest path analysis
	fmt.Printf("[>] Performing shortest path analysis on %v vertex pairs\n",
		totalNumberOfPairs(&config.Entities.DataSources))
	t3 := time.Now()
	performBfs(graph, config.Entities, config.Output)
	fmt.Printf("[>] Shortest path analysis completed in %v\n", time.Now().Sub(t3))

	// Complete
	fmt.Printf("[>] Results located at: %v\n", config.Output.OutputFile)
	fmt.Printf("[>] Total time taken: %v\n", time.Now().Sub(t0))
}

func main() {
	println("Shortest path calculator using a bipartite to unipartite transformation and the")
	println("Breadth First Search and exhaustive search algorithms with reachable vertex optimisation step")

	PerformBfsFromConfig("./config.json")
}
