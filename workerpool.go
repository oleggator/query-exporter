package main

import (
	"sync"
)

type WorkerPool struct {
	workerFunc   func(id int, tasks <-chan interface{})
	workersCount int
	tasks        chan interface{}
	wg           sync.WaitGroup
}

func NewWorkerPool(workersCount int, workerFunc func(id int, tasks <-chan interface{})) WorkerPool {
	return WorkerPool{
		workersCount: workersCount,
		workerFunc:   workerFunc,
		tasks:        make(chan interface{}, workersCount),
	}
}

func (wp *WorkerPool) Run() {
	for i := 0; i < wp.workersCount; i++ {
		wp.wg.Add(1)

		go func(id int, tasks <-chan interface{}) {
			defer wp.wg.Done()
			wp.workerFunc(id, tasks)
		}(i, wp.tasks)
	}
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

func (wp *WorkerPool) Shutdown() {
	close(wp.tasks)
}

func (wp *WorkerPool) GetInputChannel() chan<- interface{} {
	return wp.tasks
}

func (wp *WorkerPool) GetWorkersCount() int {
	return wp.workersCount
}
