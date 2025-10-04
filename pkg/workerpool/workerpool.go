package workerpool

import (
	"sync"
)

type JobResult struct {
	Output any
	Error  error
}
type JobInput struct {
	input any
}
type Job struct {
	Execute func(input JobInput) JobResult
	input   JobInput
}
type WorkerPool struct {
	jobs    chan Job
	size    int
	wg      sync.WaitGroup
	results chan JobResult
}

func NewWorkerPool(size int, bufferSize int) *WorkerPool {
	return &WorkerPool{
		jobs:    make(chan Job, bufferSize),
		size:    size,
		wg:      sync.WaitGroup{},
		results: make(chan JobResult, bufferSize),
	}
}

func (wp *WorkerPool) AddJob(job Job) {
	wp.jobs <- job
}

func (wp *WorkerPool) Start() {
	for range wp.size {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for job := range wp.jobs {
				wp.results <- job.Execute(job.input)
			}
		}()
	}
}

func (wp *WorkerPool) Close() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}
