package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	mapreduce "./mapReduce"
	"./models"
)

/* Variables globales */
var beta float64
var initialPageRank float64
var pageRanks map[int]float64
var nodes map[int]float64
var wg sync.WaitGroup

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
	cont := 0
	for scanner.Scan() {
		if cont == 0 {
			if s, err := strconv.ParseFloat(scanner.Text(), 64); err == nil {
				beta = s
				fmt.Println("Beta is ", beta)
			}
			check(err)
			break
		}
	}
	if err := scanner.Err(); err != nil {
		check(err)
	}
}

func find_lines(fileName string) chan models.Line {
	output := make(chan models.Line)

	go func() {
		defer wg.Done()
		file, err := os.Open("./inputFile/testFile.txt")
		check(err)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		cont := 0
		for scanner.Scan() {
			if cont == 0 || cont == 1 {
				continue
			}
			line := models.Line{Id: (cont - 1)}
			outgoing := strings.Split(scanner.Text(), " ")
			outgoingTotal := len(outgoing)
			line.Out = make([]int, outgoingTotal, outgoingTotal)
			for i := range outgoing {
				if s, err := strconv.Atoi(outgoing[i]); err == nil {
					line.Out = append(line.Out, s)
				}
				check(err)
			}
			output <- line
		}
		if err := scanner.Err(); err != nil {
			check(err)
		}
		//defer wg.Wait()
		close(output)
	}()
	return output
}

func mapFunc(line models.Line, output chan interface{}) {
	fmt.Println("Enter Mapper")
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
		var edge = models.Edge{Src_id: vertex.Id, Dest_id: line.Out[i], PageRank: outGoingPageRank}
		vertex.Edges = append(vertex.Edges, edge)
	}
	fmt.Println(vertex)
	results[line.Id] = vertex
	output <- results
}

func reducer(input chan interface{}, output chan interface{}) {
	results := map[string]int{}

	// for new_matches := range input {
	// 	for key, value := range new_matches.(map[string]int) {
	// 		previous_count, exists := results[key]

	// 		if !exists {
	// 			results[key] = value
	// 		} else {
	// 			results[key] = previous_count + value
	// 		}
	// 	}
	// }

	output <- results
}

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Procesando.....")
	readBetha("./inputFile/testFile.txt")
	initialPageRank = 0.2
	//input = readLines()
	//nodes := NodeCollection{m: make(map[int]*models.Node)}
	mapreduce.MapReduce(mapFunc, reducer, find_lines("./inputFile/testFile.txt"), wg)
}
