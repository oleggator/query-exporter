package main

import (
	"sync"
)

// WorkerPool manages goroutine pool
type WorkerPool struct {
	workerFunc   func(id int, tasks <-chan interface{})
	workersCount int
	tasks        chan interface{}
	wg           sync.WaitGroup
}

// NewWorkerPool creates WorkerPool
func NewWorkerPool(workersCount int, workerFunc func(id int, tasks <-chan interface{})) WorkerPool {
	return WorkerPool{
		workersCount: workersCount,
		workerFunc:   workerFunc,
		tasks:        make(chan interface{}, workersCount),
	}
}

// Run starts goroutines and tasks handling
func (wp *WorkerPool) Run() {
	for i := 0; i < wp.workersCount; i++ {
		wp.wg.Add(1)

		go func(id int, tasks <-chan interface{}) {
			defer wp.wg.Done()
			wp.workerFunc(id, tasks)
		}(i, wp.tasks)
	}
}

// Wait waits until all goroutines stopped
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// Shutdown closes channel
func (wp *WorkerPool) Shutdown() {
	close(wp.tasks)
}

// GetInputChannel returns tasks input channel
func (wp *WorkerPool) GetInputChannel() chan<- interface{} {
	return wp.tasks
}

// GetWorkersCount returns workers count
func (wp *WorkerPool) GetWorkersCount() int {
	return wp.workersCount
}
