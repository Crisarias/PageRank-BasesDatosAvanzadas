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
var inComingLinks map[int]*models.Incomings
var convergence bool

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readBetha(filename string) {
	file, err := os.Open(filename)
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if s, err := strconv.ParseFloat(scanner.Text(), 64); err == nil {
			beta = s
			fmt.Println("Beta is ", beta)
		}
		check(err)
		break
	}
	if err := scanner.Err(); err != nil {
		check(err)
	}
}

func readNodesCount(filename string) {
	file, err := os.Open(filename)
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	cont := 0
	for scanner.Scan() {
		if cont == 1 {
			if s, err := strconv.Atoi(scanner.Text()); err == nil {
				nodesCount = s
				fmt.Println("Node count is ", beta)
			}
			check(err)
			break
		}
		cont++
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
			if cont == 0 || cont == 1 {
				cont++
				continue
			}
			line := models.Line{Id: (cont - 1)}
			outgoing := strings.Split(scanner.Text(), " ")
			outgoingTotal := len(outgoing)
			line.Out = make([]int, outgoingTotal, outgoingTotal)
			for i := range outgoing {
				if s, err := strconv.Atoi(outgoing[i]); err == nil {
					line.Out[i] = s
				}
				check(err)
			}
			output <- line
			cont++
		}
		if err := scanner.Err(); err != nil {
			check(err)
		}
		close(output)
	}()
	return output
}

func aggregation_input() chan models.Incomings {
	output := make(chan models.Incomings)

	go func() {
		for _, vertex := range inComingLinks {
			output <- *vertex
		}
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
		vertex.PageRank = initialPageRank
	}
	outgoingTotal := len(line.Out)
	outGoingPageRank := vertex.PageRank / float64(outgoingTotal)
	vertex.Edges = make([]models.Edge, outgoingTotal, outgoingTotal)
	for i := range line.Out {
		vertex.Edges[i] = models.Edge{Src_id: vertex.Id, Dest_id: line.Out[i], PageRank: outGoingPageRank}
	}
	results[vertex.Id] = vertex
	output <- results
}

func mapFuncAggr(vertex models.Incomings, output chan interface{}) {
	results := map[int]models.Incomings{}
	results[vertex.Id] = vertex
	output <- results
}

func reducer(input chan interface{}, output chan interface{}) {
	results := map[int]models.Vertex{}

	for new_matches := range input {
		for _, vertex := range new_matches.(map[int]models.Vertex) {
			for _, edge := range vertex.Edges {
				if _, ok := inComingLinks[edge.Dest_id]; ok {
					inComingLinks[edge.Dest_id].SumPageRanks = inComingLinks[edge.Dest_id].SumPageRanks + edge.PageRank
				} else {
					inComingLinks[edge.Dest_id] = &models.Incomings{Id: edge.Dest_id}
					inComingLinks[edge.Dest_id].SumPageRanks = inComingLinks[edge.Dest_id].SumPageRanks + edge.PageRank
				}
			}
			if _, ok := inComingLinks[vertex.Id]; ok {
				inComingLinks[vertex.Id].Outgoings = vertex.Edges
			} else {
				inComingLinks[vertex.Id] = &models.Incomings{Id: vertex.Id}
				inComingLinks[vertex.Id].Outgoings = vertex.Edges
			}
		}
	}
	output <- results
}

func reducerAggr(input chan interface{}, output chan interface{}) {
	results := map[int]models.Incomings{}

	for new_matches := range input {
		for _, vertex := range new_matches.(map[int]models.Incomings) {
			var pageRank = ((1.0 - beta) / float64(nodesCount)) + (beta * vertex.SumPageRanks) // Change for Formula
			oldPageRanks[vertex.Id] = pageRanks[vertex.Id]
			pageRanks[vertex.Id] = pageRank
		}
	}
	output <- results
}

func validate() {
	convergence = true
	for index := 0; index < nodesCount; index++ {
		if math.Abs(pageRanks[index]-oldPageRanks[index]) > 0.0001 {
			convergence = false
			break
		}
	}

}

func main() {
	fmt.Println("Processing....")
	readBetha("./inputFile/testFile.txt")
	readNodesCount("./inputFile/testFile.txt")
	initialPageRank = 1
	pageRanks = make(map[int]float64, nodesCount)
	convergence = false
	//Start iteration cycle
	cont := 0
	//Iterate until convergence or iteration 20
	for !convergence && cont < 1000 {
		inComingLinks = make(map[int]*models.Incomings, nodesCount)
		oldPageRanks = make(map[int]float64, nodesCount)
		var wg sync.WaitGroup
		fmt.Println("Calculating Edges.....")
		// Get all page ranks for edges
		mapreduce.MapReduce(mapFunc, reducer, find_lines("./inputFile/testFile.txt"), &wg, 20)
		wg.Wait() //Wait all reducers to finish in order to have all the edges with all the information
		for _, value := range inComingLinks {
			fmt.Println(value)
		}
		var wg2 sync.WaitGroup
		fmt.Println("Calculating Page Ranks.....")
		//Calculate all page ranks for vertex based on sum of the page ranks of the incoming links edges
		mapreduce.MapReduceAggregation(mapFuncAggr, reducerAggr, aggregation_input(), &wg2, 20)
		//Wait all reducers to finish in order to have all the vertex with all the information
		wg2.Wait()
		//Validate threshold
		validate()
		fmt.Println("Iteration", cont)
		fmt.Println(pageRanks)
		cont++
	}

}
