package main

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	actual := readConfig("./test-data/test-config.json")
	expected := PathConfig{
		InputFiles: []string{"./test-data/entity_1.csv", "./test-data/entity_2.csv", "./test-data/entity_3.csv"},
		Entities: EntityConfig{
			DataSources: []DataSource{
				{
					Name:      "set-1",
					EntityIds: []string{"e-1", "e-2"},
				},
				{
					Name:      "set-2",
					EntityIds: []string{"e-5", "e-6"},
				},
			},
			Skip: []string{},
		},
		Output: OutputConfig{
			MaxDepth:        3,
			OutputFile:      "./test-data/results.csv",
			OutputDelimiter: ",",
			PathDelimiter:   "|",
			WebAppLink:      "http://192.168.99.100:8080/show/<ENTITY_IDS>",
		},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}
}

func TestNewPathResult(t *testing.T) {
	actual := NewPathResult("e-1", "e-3", []string{"e-1", "e-20", "e-3"}, "http://localhost/show.php?<ENTITY_IDS>&v")

	expected := PathResult{
		SourceEntityID:      "e-1",
		DestinationEntityID: "e-3",
		NumberOfHops:        2,
		Path:                []string{"e-1", "e-20", "e-3"},
		WebAppLink:          "http://localhost/show.php?e-1,e-20,e-3&v",
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}
}

func TestPathResultDisplay(t *testing.T) {
	pathResult := NewPathResult("e-1", "e-3", []string{"e-1", "e-20", "e-3"}, "http://localhost/show.php?<ENTITY_IDS>&v")
	actual := pathResult.display()
	expected := "e-1 -> e-3 (2 hops): [e-1 e-20 e-3]"

	if expected != actual {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}
}

func TestPathResultToString(t *testing.T) {
	pathResult := NewPathResult("e-1", "e-3", []string{"e-1", "e-20", "e-3"}, "http://localhost/show.php?<ENTITY_IDS>&v")
	actual := pathResult.toString(",", "|")
	expected := "e-1,e-3,2,e-1|e-20|e-3,http://localhost/show.php?e-1,e-20,e-3&v"

	if expected != actual {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}
}

func TestPathResultHeader(t *testing.T) {
	actual := pathResultHeader(",")
	expected := "Source entity ID,Destination entity ID,Number of hops,Path,Link"

	if expected != actual {
		t.Fatalf("Expected %v, got %v\n", expected, actual)
	}
}

func TestExtractEntityPairValid(t *testing.T) {
	e1, e2, err := extractEntityPair("e-1|e-4", "|")

	if err != nil {
		t.Fatalf("Didn't expect an error, got: %v\n", err)
	}

	if e1 != "e-1" && e2 != "e-4" {
		t.Fatalf("Entities are not as expected\n")
	}
}

func TestExtractEntityPairInvalid1(t *testing.T) {
	_, _, err := extractEntityPair("e-1", "|")

	if err == nil {
		t.Fatalf("Expected an error\n")
	}
}

func TestExtractEntityPairInvalid2(t *testing.T) {
	_, _, err := extractEntityPair("e-1|e-3|e-4", "|")

	if err == nil {
		t.Fatalf("Expected an error\n")
	}
}

func TestBuildWebAppLink(t *testing.T) {

	template := "http://192.168.99.100:8080/show.php?<ENTITY_IDS>&v"
	entityIds := []string{"e-1", "e-2"}
	actual := buildWebAppLink(template, entityIds)
	expected := "http://192.168.99.100:8080/show.php?e-1,e-2&v"

	if expected != actual {
		t.Fatalf("Expected URL: %v, got %v\n", expected, actual)
	}
}

func TestPerformBfs(t *testing.T) {

	// Build a graph with 19 vertices
	graph := NewGraph()

	// First connected component
	graph.AddUndirected("e-1", "e-2")

	// Second connected component
	graph.AddUndirected("e-3", "e-4")
	graph.AddUndirected("e-4", "e-5")
	graph.AddUndirected("e-4", "e-6")

	graph.AddUndirected("e-3", "e-8")
	graph.AddUndirected("e-8", "e-11")
	graph.AddUndirected("e-3", "e-9")
	graph.AddUndirected("e-9", "e-11")
	graph.AddUndirected("e-11", "e-13")

	graph.AddUndirected("e-3", "e-7")
	graph.AddUndirected("e-7", "e-10")
	graph.AddUndirected("e-10", "e-12")

	graph.AddUndirected("e-3", "e-14")
	graph.AddUndirected("e-3", "e-15")
	graph.AddUndirected("e-3", "e-16")
	graph.AddUndirected("e-14", "e-17")
	graph.AddUndirected("e-15", "e-17")
	graph.AddUndirected("e-16", "e-17")
	graph.AddUndirected("e-17", "e-18")
	graph.AddUndirected("e-18", "e-19")

	// Define entity pairs config
	entityConfig := EntityConfig{
		DataSources: []DataSource{
			{
				Name: "set-1",
				EntityIds: []string{
					"e-1",
					"e-2",
					"e-3",
					"e-6",
					"e-8"},
			},
			{
				Name: "set-2",
				EntityIds: []string{
					"e-11",
					"e-12",
					"e-13",
					"e-15",
					"e-16",
					"e-17",
					"e-18",
					"e-19",
					"e-100"},
			},
		},
		Skip: []string{},
	}

	// Define output config
	outputConfig := OutputConfig{
		MaxDepth:        3,
		OutputFile:      "./test-data/results.csv",
		OutputDelimiter: ",",
		PathDelimiter:   "|",
		WebAppLink:      "http://192.168.99.100:8080/show/<ENTITY_IDS>",
	}

	// Run BFS
	performBfs(&graph, entityConfig, outputConfig)

	// Check the result
	actual, err := ioutil.ReadFile("./test-data/results.csv")
	if err != nil {
		t.Fatalf("Unable to find test results\n")
	}

	expected, err := ioutil.ReadFile("./test-data/expected_results.csv")
	if err != nil {
		t.Fatalf("Unable to find expected results\n")
	}

	if !bytes.Equal(expected, actual) {
		t.Fatalf("Actual results differ from expected results\n")
	}

}

func TestPerformBfsFromConfig(t *testing.T) {

	// Perform BFS using bipartite data
	PerformBfsFromConfig("./test-data-full/config.json")

	// Check the result
	actual, err := ioutil.ReadFile("./test-data-full/results.csv")
	if err != nil {
		t.Fatalf("Unable to find test results\n")
	}

	expected, err := ioutil.ReadFile("./test-data-full/expected_results.csv")
	if err != nil {
		t.Fatalf("Unable to find expected results\n")
	}

	if !bytes.Equal(expected, actual) {
		t.Fatalf("Actual results differ from expected results\n")
	}
}

func TestPerformBfsFromConfigThreeDataSources(t *testing.T) {

	// Perform BFS using bipartite data
	PerformBfsFromConfig("./test-data-full/config-2.json")

	// Check the result
	actual, err := ioutil.ReadFile("./test-data-full/results-2.csv")
	if err != nil {
		t.Fatalf("Unable to find test results\n")
	}

	expected, err := ioutil.ReadFile("./test-data-full/expected_results-2.csv")
	if err != nil {
		t.Fatalf("Unable to find expected results\n")
	}

	if !bytes.Equal(expected, actual) {
		t.Fatalf("Actual results differ from expected results\n")
	}
}

func TestPerformBfsFromConfigWithSkips(t *testing.T) {

	// Perform BFS using bipartite data
	PerformBfsFromConfig("./test-data-full-2/config.json")

	// Check the result
	actual, err := ioutil.ReadFile("./test-data-full-2/results.csv")
	if err != nil {
		t.Fatalf("Unable to find test results\n")
	}

	expected, err := ioutil.ReadFile("./test-data-full-2/expected_results.csv")
	if err != nil {
		t.Fatalf("Unable to find expected results\n")
	}

	if !bytes.Equal(expected, actual) {
		t.Fatalf("Actual results differ from expected results\n")
	}
}

func TestPerformFindAllShortestPathsFromConfig(t *testing.T) {

	// Perform BFS using bipartite data
	PerformBfsFromConfig("./test-data-full-3/config.json")

	// Check the result
	actual, err := ioutil.ReadFile("./test-data-full-3/results.csv")
	if err != nil {
		t.Fatalf("Unable to find test results\n")
	}

	expected, err := ioutil.ReadFile("./test-data-full-3/expected_results.csv")
	if err != nil {
		t.Fatalf("Unable to find expected results\n")
	}

	if !bytes.Equal(expected, actual) {
		t.Fatalf("Actual results differ from expected results\n")
	}
}

func TestTotalNumberOfPairsOneDataset(t *testing.T) {
	set := []DataSource{
		{
			Name:      "set-1",
			EntityIds: []string{"e-1"},
		},
	}

	actual := totalNumberOfPairs(&set)

	if actual != 0 {
		t.Errorf("Expected 0 pairs, got %v\n", actual)
	}
}

func TestTotalNumberOfPairsTwoDatasets(t *testing.T) {
	set := []DataSource{
		{
			Name:      "set-1",
			EntityIds: []string{"e-1", "e-2", "e-3"},
		},
		{
			Name:      "set-2",
			EntityIds: []string{"e-1", "e-3"},
		},
	}

	actual := totalNumberOfPairs(&set)

	if actual != 6 {
		t.Errorf("Expected 6 pairs, got %v\n", actual)
	}
}

func TestTotalNumberOfPairsThreeDatasets(t *testing.T) {
	set := []DataSource{
		{
			Name:      "set-1",
			EntityIds: []string{"e-1", "e-2", "e-3"},
		},
		{
			Name:      "set-2",
			EntityIds: []string{"e-1", "e-3"},
		},
		{
			Name:      "set-3",
			EntityIds: []string{"e-1"},
		},
	}

	actual := totalNumberOfPairs(&set)

	// set-1 and set-2 = 6
	// set-1 and set-3 = 3
	// set-2 and set-3 = 2
	// therefore expected total = 11

	if actual != 11 {
		t.Errorf("Expected 11 pairs, got %v\n", actual)
	}
}
