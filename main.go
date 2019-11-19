package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	mapreduce "./mapReduce"
	"./models"
)

/* Variables globales */
var beta float64
var nodesCount int
var initialPageRank float64

var pageRanks map[int]float64
var mutex = &sync.RWMutex{}
var inLinks map[int]*models.InLinks
var lines []models.Line

var outLinks map[int][]int
var convergence bool

const convergenceDifference = 0.00000001

const inputFile = "./inputFile/testFile.txt"

func check(e error) {
	if e != nil {
		fmt.Println("Unexpected error: %s", e)
		panic(e)
	}
}

func readNodesCount() {
	file, err := os.Open(inputFile)
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if s, err := strconv.ParseFloat(scanner.Text(), 64); err == nil {
			nodesCount = int(s)
			initialPageRank = (1 / s)
			fmt.Println("Node count is ", nodesCount)
			fmt.Println("Initial pageRank is ", initialPageRank)
		}
		check(err)
		break
	}
	if err := scanner.Err(); err != nil {
		check(err)
	}
}

func find_lines() chan models.Line {
	output := make(chan models.Line)

	go func() {
		file, err := os.Open(inputFile)
		check(err)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		cont := 0
		for scanner.Scan() {
			if cont > 0 {
				line := models.Line{Id: (cont - 1)}
				outgoing := strings.Split(scanner.Text(), " ")
				outgoingTotal := len(outgoing)
				line.Out = make([]int, outgoingTotal, outgoingTotal)
				outLinks[cont] = make([]int, outgoingTotal, outgoingTotal)
				for i := range outgoing {
					if s, err := strconv.Atoi(outgoing[i]); err == nil {
						line.Out[i] = s
						outLinks[cont][i] = s
					}
					check(err)
				}
				lines[cont-1] = line
				output <- line
			}
			cont++
		}
		if err := scanner.Err(); err != nil {
			check(err)
		}
		close(output)
	}()
	return output
}

func send_lines() chan models.Line {
	output := make(chan models.Line)
	go func() {
		for i := 0; i < nodesCount; i++ {
			output <- lines[i]
		}
		close(output)
	}()
	return output
}

func aggregation_input() chan models.InLinks {
	output := make(chan models.InLinks)
	cont := nodesCount
	go func() {
		for i := 0; i < cont; i++ {
			if _, ok := inLinks[i]; ok {
				output <- *inLinks[i]
			} else {
				//If the node has no inlinks pass to the reducer with a sumPageRanks of 0
				inLinks[i] = &models.InLinks{Id: i}
				inLinks[i].SumPageRanks = 0
				output <- *inLinks[i]
			}
		}
		close(output)
	}()
	return output
}

func mapFunc(line models.Line, output chan interface{}) {
	results := map[int]models.Vertex{}
	var vertex = models.Vertex{Id: line.Id}
	//Get PageRank if not exist assign initial
	mutex.RLock()
	if val, ok := pageRanks[vertex.Id]; ok {
		mutex.RUnlock()
		vertex.PageRank = val
	} else {
		mutex.RUnlock()
		mutex.Lock()
		pageRanks[vertex.Id] = initialPageRank
		mutex.Unlock()
		vertex.PageRank = initialPageRank
	}
	//If have no outlinks has to link to all
	if line.Out[0] == 0 {
		outGoingPageRank := vertex.PageRank / float64(nodesCount-1)
		vertex.Edges = make([]models.Edge, nodesCount-1, nodesCount-1)
		index := 0
		for i := 0; i < nodesCount; i++ {
			if (i) != vertex.Id {
				vertex.Edges[index] = models.Edge{Dest_id: i, PageRank: outGoingPageRank}
				index++
			}
		}
	} else {
		outgoingTotal := len(line.Out)
		outGoingPageRank := vertex.PageRank / float64(outgoingTotal)
		vertex.Edges = make([]models.Edge, outgoingTotal, outgoingTotal)
		for i := range line.Out {
			vertex.Edges[i] = models.Edge{Dest_id: line.Out[i], PageRank: outGoingPageRank}
		}
	}

	results[vertex.Id] = vertex
	output <- results
}

func mapFuncAggregation(vertex models.InLinks, output chan interface{}) {
	results := map[int]models.InLinks{}
	results[vertex.Id] = vertex
	output <- results
}

func reducer(input chan interface{}, output chan interface{}) {
	results := map[int]models.Vertex{}

	for new_matches := range input {
		for _, vertex := range new_matches.(map[int]models.Vertex) {
			for _, edge := range vertex.Edges {
				if _, ok := inLinks[edge.Dest_id]; ok {
					inLinks[edge.Dest_id].SumPageRanks = inLinks[edge.Dest_id].SumPageRanks + edge.PageRank
				} else {
					inLinks[edge.Dest_id] = &models.InLinks{Id: edge.Dest_id}
					inLinks[edge.Dest_id].SumPageRanks = inLinks[edge.Dest_id].SumPageRanks + edge.PageRank
				}
			}
		}
	}
	output <- results
}

func reducerAggregation(input chan interface{}, output chan interface{}) {
	results := map[int]models.InLinks{}
	for new_matches := range input {
		for _, vertex := range new_matches.(map[int]models.InLinks) {
			var pageRank = ((1.0 - beta) / float64(nodesCount)) + (beta * (vertex.SumPageRanks))
			mutex.RLock()
			old, ok := pageRanks[vertex.Id]
			mutex.RUnlock()
			mutex.Lock()
			pageRanks[vertex.Id] = pageRank
			mutex.Unlock()
			if convergence && ok && (math.Abs(pageRank-old) > convergenceDifference) {
				convergence = false
			}
		}
	}
	output <- results
}

//Write results files
func endProcess(iteration int) {
	totalsum := 0.0
	fmt.Println("Writing final pageRanks")
	os.Remove("Results.txt")
	file, err := os.OpenFile("Results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err)
	datawriter := bufio.NewWriter(file)
	for i := 0; i < nodesCount; i++ {
		pagerank, _ := pageRanks[i]
		totalsum += pagerank
		_, _ = datawriter.WriteString(strconv.FormatFloat(pagerank, 'f', 25, 64) + "\n")
		//fmt.Println("Edge ", i, ", ", "Page Rank: ", pagerank.(float64), ", Outlinks:", outLinks[i])
		fmt.Println("Edge ", i, ", ", "Page Rank: ", pagerank)
	}
	datawriter.Flush()
	file.Close()
	fmt.Println("All nodes converge with +/- ", convergenceDifference, " after iteration", iteration)
	fmt.Println("Sum of all page ranks ", totalsum)
	fmt.Println("InitialPageRank", initialPageRank)
}

func main() {
	fmt.Println("Running on ", runtime.NumCPU(), " cores")
	runtime.GOMAXPROCS(runtime.NumCPU())
	beta = -1
	for beta == -1 {
		fmt.Println("Enter a value for betha")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if s, err := strconv.ParseFloat(scanner.Text(), 64); err == nil {
				if s <= 0 || s >= 1 {
					fmt.Println("Betha value must higher than 0 and lower than 1")
				} else {
					beta = s
					fmt.Println("Betha from terminal is ", beta)
					break
				}
			} else {
				fmt.Println("Please a valid float64 number")
			}
		}
		if scanner.Err() != nil {
			fmt.Println("There was an error reading from the terminal")
		}
	}
	fmt.Println("Processing....")
	readNodesCount()
	lines = make([]models.Line, nodesCount)
	pageRanks = make(map[int]float64, nodesCount)
	outLinks = make(map[int][]int, nodesCount)
	convergence = false
	//Start iteration cycle
	cont := 0
	//Set pool size
	poolSize := float64(nodesCount) * 0.0002
	//Iterate until convergence
	for !convergence {
		convergence = true
		inLinks = make(map[int]*models.InLinks, nodesCount)
		var wg sync.WaitGroup
		//fmt.Println("Calculating Edges.....")
		// Get all page ranks for edges
		if cont == 0 {
			mapreduce.MapReduce(mapFunc, reducer, find_lines(), &wg, int(poolSize))
		} else {
			mapreduce.MapReduce(mapFunc, reducer, send_lines(), &wg, int(poolSize))
		}
		wg.Wait() //Wait all reducers to finish in order to have all the edges with all the information
		var wg2 sync.WaitGroup
		//fmt.Println("Calculating Page Ranks.....")
		//Calculate all page ranks for vertex based on sum of the page ranks of the incoming links edges
		mapreduce.MapReduceAggregation(mapFuncAggregation, reducerAggregation, aggregation_input(), &wg2, int(poolSize))
		//Wait all reducers to finish in order to have all the vertex with all the information
		wg2.Wait()
		fmt.Println("Iteration", cont, "Completed")
		cont++
	}
	endProcess(cont)
}
