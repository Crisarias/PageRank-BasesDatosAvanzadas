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

func line(filename string) chan string {
	output := make(chan string)

	go func() {
		file, err := os.Open(filename)
		check(err)
		defer file.Close()
		cont := 0
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			if cont == 0 || cont == 1 {
				cont++
				continue
			}
			s := strconv.Itoa(cont - 1)
			newLine := s + " " + scanner.Text()
			fmt.Println(newLine)
			output <- newLine
			cont++
		}
		if err := scanner.Err(); err != nil {
			check(err)
		}

		close(output)
	}()

	return output
}

func find_files(dirname string) chan interface{} {
	output := make(chan interface{})

	go func() {
		_find_files(dirname, output)
		close(output)
	}()

	return output
}

func _find_files(dirname string, output chan interface{}) {
	dir, _ := os.Open(dirname)
	dirnames, _ := dir.Readdirnames(-1)

	for i := 0; i < len(dirnames); i++ {
		fullpath := dirname + "/" + dirnames[i]
		file, _ := os.Stat(fullpath)

		if file.IsDir() {
			_find_files(fullpath, output)
		} else {
			output <- fullpath
		}
	}
}

func mapFunc(filename interface{}, output chan interface{}) {
	fmt.Println("Enter Mapper")
	results := map[int]models.Vertex{}

	for line := range line(filename.(string)) {
		outgoing := strings.Split(line, " ")
		outgoingTotal := len(outgoing) - 1
		var vertex models.Vertex
		if n, err := strconv.Atoi(outgoing[0]); err == nil {
			vertex.Id = n
		}
		//ObtenerPageRank
		if val, ok := pageRanks[vertex.Id]; ok {
			vertex.PageRank = val
		} else {
			vertex.PageRank = initialPageRank
		}
		outGoingPageRank := vertex.PageRank / float64(outgoingTotal)
		vertex.Edges = make([]models.Edge, outgoingTotal, outgoingTotal)
		for i := range outgoing {
			if i == 0 {
				continue
			}
			if n, err := strconv.Atoi(outgoing[i]); err == nil {
				fmt.Println(outGoingPageRank)
				vertex.Edges[i-1] = models.Edge{Src_id: vertex.Id, Dest_id: n, PageRank: outGoingPageRank}
			}
		}
		fmt.Println(vertex)
		results[vertex.Id] = vertex
	}

	output <- results
}

func reducer(input chan interface{}, output chan interface{}) {
	fmt.Println("Reducer")
	results := map[int]models.Vertex{}

	for new_matches := range input {
		for key, value := range new_matches.(map[int]models.Vertex) {
			previous_count, exists := results[key]
			fmt.Println("Reducer loop")
			fmt.Println(value)
			if !exists {
				results[key] = value
			} else {
				fmt.Println(previous_count)
				results[key] = value
			}
		}
	}

	output <- results
}

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Procesando.....")
	//readBetha("./inputFiles/testFile.txt")
	initialPageRank = 0.2
	//input = readLines()
	//nodes := NodeCollection{m: make(map[int]*models.Node)}
	mapreduce.MapReduce(mapFunc, reducer, find_files("./inputFiles"), 20)
}
