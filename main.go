package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
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
var oldPageRanks map[int]float64
var inLinks map[int]*models.InLinks
var outLinks map[int][]int
var convergence bool

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readNodesCount(filename string) {
	file, err := os.Open(filename)
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if s, err := strconv.Atoi(scanner.Text()); err == nil {
			nodesCount = s
			fmt.Println("Node count is ", nodesCount)
		}
		check(err)
		break
	}
	if err := scanner.Err(); err != nil {
		check(err)
	}
}

func find_lines(fileName string) chan models.Line {
	output := make(chan models.Line)

	go func() {
		file, err := os.Open("./inputFile/testFile.txt")
		check(err)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		cont := 0
		for scanner.Scan() {
			if cont > 0 {
				line := models.Line{Id: (cont)}
				outgoing := strings.Split(scanner.Text(), " ")
				// outgoingTotal := 0
				// if outgoing[0] != " " {
				// 	outgoingTotal = len(outgoing)
				// }
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

func aggregation_input() chan models.InLinks {
	output := make(chan models.InLinks)
	cont := nodesCount
	go func() {
		for i := 1; i <= cont; i++ {
			if _, ok := inLinks[i]; ok {
				output <- *inLinks[i]
			} else {
				fmt.Println("No inlinks")
				inLinks[i] = &models.InLinks{Id: i}
				inLinks[i].SumPageRanks = 0
				output <- *inLinks[i]
			}
		}
		// for _, vertex := range inLinks {
		// 	output <- *vertex
		// }
		close(output)
	}()
	return output
}

func mapFunc(line models.Line, output chan interface{}) {
	results := map[int]models.Vertex{}
	var vertex = models.Vertex{Id: line.Id}
	//ObtenerPageRank
	if val, ok := pageRanks[vertex.Id]; ok {
		vertex.PageRank = val
	} else {
		pageRanks[vertex.Id] = initialPageRank
		vertex.PageRank = initialPageRank
	}
	//If have no outlinks has to link to all
	if line.Out[0] == 0 {
		outGoingPageRank := vertex.PageRank / float64(nodesCount-1)
		vertex.Edges = make([]models.Edge, nodesCount-1, nodesCount-1)
		index := 0
		for i := 1; i <= nodesCount; i++ {
			if (i) != vertex.Id {
				fmt.Println("Added")
				vertex.Edges[index] = models.Edge{Src_id: vertex.Id, Dest_id: i, PageRank: outGoingPageRank}
				index++
			}
		}
	} else {
		outgoingTotal := len(line.Out)
		outGoingPageRank := vertex.PageRank / float64(outgoingTotal)
		vertex.Edges = make([]models.Edge, outgoingTotal, outgoingTotal)
		for i := range line.Out {
			vertex.Edges[i] = models.Edge{Src_id: vertex.Id, Dest_id: line.Out[i], PageRank: outGoingPageRank}
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
		fmt.Println("Enter aggregation")
		for _, vertex := range new_matches.(map[int]models.InLinks) {
			var pageRank = ((1.0 - beta) / float64(nodesCount)) + (beta * (vertex.SumPageRanks))
			oldPageRanks[vertex.Id] = pageRanks[vertex.Id]
			pageRanks[vertex.Id] = pageRank
		}
	}
	output <- results
}

func validate() {
	convergence = true
	for index := 0; index < nodesCount; index++ {
		if math.Abs(pageRanks[index]-oldPageRanks[index]) != 0 {
			convergence = false
			break
		}
	}
}

func main() {
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
	readNodesCount("./inputFile/testFile.txt")
	initialPageRank = 1
	pageRanks = make(map[int]float64, nodesCount)
	outLinks = make(map[int][]int, nodesCount)
	convergence = false
	//Start iteration cycle
	cont := 0
	//Iterate until convergence or iteration 20
	for !convergence && cont < 1000 {
		inLinks = make(map[int]*models.InLinks, nodesCount)
		oldPageRanks = make(map[int]float64, nodesCount)
		var wg sync.WaitGroup
		fmt.Println("Calculating Edges.....")
		// Get all page ranks for edges
		mapreduce.MapReduce(mapFunc, reducer, find_lines("./inputFile/testFile.txt"), &wg, 20)
		wg.Wait() //Wait all reducers to finish in order to have all the edges with all the information
		var wg2 sync.WaitGroup
		fmt.Println("Calculating Page Ranks.....")
		//Calculate all page ranks for vertex based on sum of the page ranks of the incoming links edges
		mapreduce.MapReduceAggregation(mapFuncAggregation, reducerAggregation, aggregation_input(), &wg2, 20)
		//Wait all reducers to finish in order to have all the vertex with all the information
		wg2.Wait()
		//Validate threshold
		validate()
		fmt.Println("Iteration", cont)
		fmt.Println(pageRanks)
		cont++
	}
	totalsum := 0.0
	fmt.Println("Final vector with outgoings")
	for i := 1; i <= nodesCount; i++ {
		totalsum += pageRanks[i]
		fmt.Println("Edge ", i, ", ", "Page Rank: ", pageRanks[i], ", Outlinks:", outLinks[i])
	}
	fmt.Println("Final Total", totalsum)

}
