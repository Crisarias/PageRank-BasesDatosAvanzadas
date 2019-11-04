package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"

	"./models"
)

/* Variables globales */
var beta float64
var initialPageRank float64
var pageRanks map[int]float64

type NodeCollection struct {
	sync.RWMutex
	m map[int]*models.Node
}

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
				break;
			}
		}
		if err := scanner.Err(); err != nil {
			check(err)
		}
}

func mapStepPropagation(filename string) chan  models.Vertex {
	output_vertex := make(chan models.Vertex)
	go func() {
		file, err := os.Open(filename)
		check(err)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		cont := 0
		for scanner.Scan() {
			if cont == 0 || cont == 1 {
				continue
			}
			vertex = models.Vertex{id:cont-1}
			//ObtenerPageRank
			if val, ok := pageRanks[vertex.id]; ok {
				vertex.pageRank = val
			}
			else {
				vertex.pageRank = initialPageRank
			}			
			outgoing := strings.Split(scanner.Text(), " ")
			outgoingTotal := len(outgoing.count)
			outGoingPageRank := vertex.pageRank / outgoingTotal
			vertex.edges := make([]models.Edge, outgoingTotal, outgoingTotal) 
			for _, var := range  outgoing{
				if s, err := strconv.ParseInt(var, 64); err == nil {
					edge = models.Edge{src_id:vertex.id, dest_id:s, pageRank: outGoingPageRank}
					vertex.edges = append(vertex.edges, edge)
				}
				check(err)				
			}
			// Agregar la linea al canal add each json line to our enumeration channel
			output_vertex <- vertex
			fmt.Println(vertex)
			cont++
		}
		if err := scanner.Err(); err != nil {
			check(err)
		}
		close(output)
	}()
	return output
}


func mapStepAggregation(filename string) chan  models.Vertex, chan models.Edge  {
	output_vertex := make(chan models.Vertex)
	output_edges := make(chan models.Edge)
	go func() {
		file, err := os.Open(filename)
		check(err)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		cont := 0
		for scanner.Scan() {
			if cont == 0 || cont == 1 {
				continue
			}
			vertex = models.Vertex{id:cont-1}
			//ObtenerPageRank
			if val, ok := pageRanks[vertex.id]; ok {
				vertex.pageRank = val
			}
			else {
				vertex.pageRank = initialPageRank
			}
			// Agregar la linea al canal add each json line to our enumeration channel
			output_vertex <- vertex
			fmt.Println(vertex)
			outgoing := strings.Split(scanner.Text(), " ")
			outgoingTotal := len(outgoing.count)
			outGoingPageRank := vertex.pageRank / outgoingTotal
			for _, var := range  outgoing{
				if s, err := strconv.ParseInt(var, 64); err == nil {
					edge = models.Edge{src_id:vertex.id, dest_id:s, pageRank: outGoingPageRank}
					output_edges <- vertex
					fmt.Println(edge)
				}
				check(err)				
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

func reducer(input chan interface{}, output chan interface{}) {
	results := map[string]int{}

	for new_matches := range input {
		for key, value := range new_matches.(map[string]int) {
			previous_count, exists := results[key]

			if !exists {
				results[key] = value
			} else {
				results[key] = previous_count + value
			}
		}
	}

	output <- results
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Procesando.....")
	readBetha()
	initialPageRank = 0.2
	input = readLines()
	//nodes := NodeCollection{m: make(map[int]*models.Node)}
}

func mapToNode(index int, text string, nodes *NodeCollection) {

}
