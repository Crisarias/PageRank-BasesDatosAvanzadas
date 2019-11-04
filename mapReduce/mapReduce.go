package mapreduce

import (
	"fmt"
	"sync"

	"../models"
)

func MapReduce(mapper func(models.Line, chan interface{}),
	reducer func(chan interface{}, chan interface{}),
	input chan models.Line, wg sync.WaitGroup) interface{} {
	fmt.Println("Enter Map function")

	reduce_input := make(chan interface{})
	reduce_output := make(chan interface{})
	//worker_output := make(chan chan interface{}, pool_size)

	go reducer(reduce_input, reduce_output)

	// go func() {
	// 	defer waitGrp.Done()
	// 	for worker_chan := range worker_output {
	// 		reduce_input <- <-worker_chan
	// 	}

	// 	close(reduce_input)
	// }()

	go func() {
		defer wg.Done()
		fmt.Println("input Is", input)
		for item := range input {
			fmt.Println(item)
			my_chan := make(chan interface{})
			go mapper(item, my_chan)
			//worker_output <- my_chan
		}
		//defer wg.Wait()
		//close(worker_output)
	}()

	return <-reduce_output
}
