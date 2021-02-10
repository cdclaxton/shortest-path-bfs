package main

import (
	"encoding/json"
	"flag"
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
	Name      string   `json:"name"`       // friendly name of the data source
	EntityIds []string `json:"entity_ids"` // list of entity IDs
}

// EntityConfig represents the entity pairs for which to find paths
type EntityConfig struct {
	DataSources []DataSource `json:"data_sources"` // list of data sources with entity IDs of interest
	Skip        []string     `json:"skip"`         // list of entities to ignore when constructing the graph
}

// OutputConfig represents the config for the output from the BFS
type OutputConfig struct {
	MaxDepth        int    `json:"max_depth"`      // maximum number of hops from a source to a destination vertex
	FindAllPaths    bool   `json:"find_all_paths"` // should all paths be found or just the first?
	OutputFile      string `json:"output_file"`    // location of the output CSV file
	OutputDelimiter string `json:"delimiter"`      // delimiter to use in the CSV file
	PathDelimiter   string `json:"path_delimiter"` // delimiter to use between entity IDs on a path
	WebAppLink      string `json:"webapp_link"`    // web-app link to generate for the path
	UnipartiteFile  string `json:"unipartite"`     // location of the unipartite CSV file to write
}

// PathConfig represents the JSON config
type PathConfig struct {
	InputFiles []string     `json:"input_files"` // list of CSV files from which the graph will be constructed
	Entities   EntityConfig `json:"entities"`    // entity IDs to consider and skip
	Output     OutputConfig `json:"output"`      // configuration for the output CSV file
}

// display the path config
func (c *PathConfig) display() {
	log.Println("Parameter - Number of input files:      ", len(c.InputFiles))
	log.Println("Parameter - Number of data sources:     ", len(c.Entities.DataSources))
	log.Println("Parameter - Number of entities to skip: ", len(c.Entities.Skip))
	log.Println("Parameter - Maximum depth:              ", c.Output.MaxDepth)
	log.Println("Parameter - Find all paths:             ", c.Output.FindAllPaths)
	log.Println("Parameter - Output file:                ", c.Output.OutputFile)
	log.Println("Parameter - Delimiter:                  ", c.Output.OutputDelimiter)
	log.Println("Parameter - Path delimiter:             ", c.Output.PathDelimiter)
	log.Println("Parameter - Web-app link template:      ", c.Output.WebAppLink)
	log.Println("Parameter - Unipartite graph file:      ", c.Output.UnipartiteFile)
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
		log.Fatalf("Unable to unmarshall JSON from file: %v", filePath)
	}

	return config
}

// PathResult represents a shortest path
type PathResult struct {
	SourceEntityID              string   // entity ID of the source vertex
	SourceEntityDataSource      string   // data source from which the source entity ID came
	DestinationEntityID         string   // entity ID of the destination vertex
	DestinationEntityDataSource string   // data source from which the destination entity ID came
	NumberOfHops                int      // number of hops from source to destination
	Path                        []string // list of entity IDs on the path from source to destination
	WebAppLink                  string   // web-app link for the path
}

// buildWebAppLink builds the web-app link
func buildWebAppLink(template string, path []string) string {

	// Precondition
	if len(path) == 0 {
		log.Fatal("Path is empty")
	}

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

	// Preconditions
	if len(delimiter) == 0 {
		log.Fatal("Cannot use a blank delimiter")
	}

	if len(pathDelimiter) == 0 {
		log.Fatal("Cannot use a blank delimiter for the path")
	}

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

	// Precondition
	if len(delimiter) == 0 {
		log.Fatal("Cannot use a blank delimiter")
	}

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

// findAndRecordShortestPaths finds the shortest path and writes to file and returns the number of paths found
func findAndRecordShortestPaths(g *Graph,
	source string, sourceDataSource string,
	destination string, destinationDataSource string,
	outputConfig OutputConfig, outputFile *os.File) int {

	numPathsFound := 0

	if outputConfig.FindAllPaths {

		// Find all the paths between the source and destination up to a maximum length
		paths := g.AllPaths(source, destination, outputConfig.MaxDepth)
		numPathsFound = len(paths)

		if len(paths) == 0 {
			log.Fatalf("Vertex %v was deemed reachable from %v, but no path!\n", destination, source)
		} else {
			for _, path := range paths {
				result := NewPathResult(source, sourceDataSource,
					destination, destinationDataSource,
					path.flatten(), outputConfig.WebAppLink)
				log.Printf("%v\n", result.display())
				fmt.Fprintln(outputFile, result.toString(outputConfig.OutputDelimiter, outputConfig.PathDelimiter))
			}
		}

	} else {
		// Compute the shortest path using BFS
		found, vertex := g.Bfs(source, destination, outputConfig.MaxDepth)

		if !found {
			log.Fatalf("Vertex %v was deemed reachable from %v, but no path!\n", destination, source)
		} else {

			// Found 1 path
			numPathsFound++

			// Build the PathResult
			result := NewPathResult(source, sourceDataSource,
				destination, destinationDataSource,
				vertex.flatten(), outputConfig.WebAppLink)

			// Display the result
			log.Printf("%v\n", result.display())

			// Add the result to the file
			fmt.Fprintln(outputFile, result.toString(outputConfig.OutputDelimiter, outputConfig.PathDelimiter))
		}
	}

	return numPathsFound
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
		log.Fatalf("Unable to open output file %v for writing: %v\n", outputConfig.OutputFile, err)
	}
	defer outputFile.Close()

	// Write the header to the output CSV file
	fmt.Fprintln(outputFile, pathResultHeader(outputConfig.OutputDelimiter))

	// Total number of entity pairs to check
	totalPairs := totalNumberOfPairs(&entityConfig.DataSources)
	numPairsProcessed := 0

	// Make a set of entities to skip
	skipEntities := SliceToSet(entityConfig.Skip)

	numPairsWithPaths := 0
	numPathsFound := 0

	// Walk through all pairs of data sources
	for i := 0; i < len(entityConfig.DataSources)-1; i++ {
		for j := i + 1; j < len(entityConfig.DataSources); j++ {

			log.Printf("Checking connections for data sources %v <--> %v\n",
				entityConfig.DataSources[i].Name,
				entityConfig.DataSources[j].Name)

			// Walk through each source entity in the i(th) dataset
			for _, source := range entityConfig.DataSources[i].EntityIds {

				// Skip the source entity if required
				if skipEntities.Has(source) {
					// Don't need to check all paths to the j(th) dataset
					numPairsProcessed += len(entityConfig.DataSources[j].EntityIds)
					continue
				}

				// Set of all vertices within reach of the source vertex
				found, reachable := g.ReachableVertices(source, outputConfig.MaxDepth)

				// If the source vertex was not found in the dataset, just continue to the next vertex
				if !found {
					numPairsProcessed += len(entityConfig.DataSources[j].EntityIds)
					continue
				}

				// Walk through each destination entity in the j(th) dataset
				for _, destination := range entityConfig.DataSources[j].EntityIds {

					// Skip the entity if it's both source and destination or if it needs to be skipped
					if (source == destination) || skipEntities.Has(destination) {
						numPairsProcessed++
						continue
					}

					// Provide feedback on long-running jobs
					if numPairsProcessed%10000 == 0 {
						log.Printf("Processed %v pairs of %v\n", numPairsProcessed, totalPairs)
					}

					// If the destination is reachable from the source, then find and record the shortest path
					if reachable.Has(destination) {
						numPathsFound += findAndRecordShortestPaths(
							g,
							source,
							entityConfig.DataSources[i].Name,
							destination,
							entityConfig.DataSources[j].Name,
							outputConfig,
							outputFile)

						numPairsWithPaths++
					}

					numPairsProcessed++
				}

			}

		}

		log.Printf("Summary - Total number of entity pairs:   %v\n", totalPairs)
		log.Printf("Summary - Number of pairs with paths:     %v\n", numPairsWithPaths)
		log.Printf("Summary - Percentage of pairs with paths: %.2f %%\n", 100.0*float32(numPairsWithPaths)/float32(totalPairs))
		log.Printf("Summary - Total number of paths found:    %v\n", numPathsFound)
	}

}

// PerformBfsFromConfig performs BFS based on a config file
func PerformBfsFromConfig(configFilepath string) {

	// Read the JSON configuration
	t0 := time.Now()
	log.Println("Reading configuration ...")
	config := readConfig(configFilepath)
	config.display()

	// Check there are at least two data sources to find connections
	if len(config.Entities.DataSources) < 2 {
		log.Println("At least two data sources must be specified in the config")
		return
	}

	// Read the entity-document relationships from file
	log.Println("Reading entity-document graph from file ...")
	t1 := time.Now()
	connections := ReadEntityDocumentGraph(config.InputFiles, SliceToSet(config.Entities.Skip))
	log.Printf("Entity-document graph read in %v\n", time.Now().Sub(t1))

	// Convert the bipartite graph to a unipartite graph
	t2 := time.Now()
	graph := BipartiteToUnipartite(connections)
	log.Printf("Graph has %v vertices\n", len(graph.Nodes))
	log.Printf("Bipartite to unipartite conversion completed in %v\n", time.Now().Sub(t2))

	// Write the unipartite graph to file (if required)
	if len(config.Output.UnipartiteFile) > 0 {
		log.Printf("Writing unipartite graph to file: %v\n", config.Output.UnipartiteFile)
		graph.WriteUndirectedEdgeList(config.Output.UnipartiteFile, config.Output.PathDelimiter)
	}

	// Perform shortest path analysis
	log.Printf("Performing shortest path analysis on %v vertex pairs\n",
		totalNumberOfPairs(&config.Entities.DataSources))
	t3 := time.Now()
	performBfs(graph, config.Entities, config.Output)
	log.Printf("Shortest path analysis completed in %v\n", time.Now().Sub(t3))

	// Complete
	log.Printf("Results located at: %v\n", config.Output.OutputFile)
	log.Printf("Total time taken: %v\n", time.Now().Sub(t0))
}

func main() {

	// Command line arguments
	configFilepath := flag.String("config", "config.json", "Location of the JSON config file")
	flag.Parse()

	log.Println("Shortest path calculator using a bipartite to unipartite transformation and the")
	log.Println("Breadth First Search and exhaustive search algorithms with reachable vertex optimisation step")

	PerformBfsFromConfig(*configFilepath)
}
