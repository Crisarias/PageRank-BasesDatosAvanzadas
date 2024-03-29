package mapreduce

import (
	"sync"

	"../models"
)

func MapReduce(mapper func(models.Line, chan interface{}),
	reducer func(chan interface{}, chan interface{}),
	input chan models.Line, wg *sync.WaitGroup, pool_size int) interface{} {
	reduce_input := make(chan interface{})
	reduce_output := make(chan interface{})
	worker_output := make(chan chan interface{}, pool_size)

	wg.Add(1)
	go func() {
		defer wg.Done()
		go reducer(reduce_input, reduce_output)

		go func() {
			for worker_chan := range worker_output {
				reduce_input <- <-worker_chan
			}

			close(reduce_input)
		}()

		go func() {
			for item := range input {
				my_chan := make(chan interface{})
				go mapper(item, my_chan)
				worker_output <- my_chan
			}
			close(worker_output)
		}()
	}()

	return <-reduce_output
}

func MapReduceAggregation(mapper func(models.InLinks, chan interface{}),
	reducer func(chan interface{}, chan interface{}),
	input chan models.InLinks, wg *sync.WaitGroup, pool_size int) interface{} {

	reduce_input := make(chan interface{})
	reduce_output := make(chan interface{})
	worker_output := make(chan chan interface{}, pool_size)

	wg.Add(1)
	go func() {
		defer wg.Done()
		go reducer(reduce_input, reduce_output)

		go func() {
			for worker_chan := range worker_output {
				reduce_input <- <-worker_chan
			}

			close(reduce_input)
		}()

		go func() {
			for item := range input {
				my_chan := make(chan interface{})
				go mapper(item, my_chan)
				worker_output <- my_chan
			}
			close(worker_output)
		}()
	}()

	return <-reduce_output
}
