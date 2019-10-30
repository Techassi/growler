package workerpool

import (
	"github.com/Techassi/growler/worker"
)

type WorkerPool struct {
	MaxWorkers        int
	ProcessingChannel chan<- bool
	FinishedChannel   chan<- worker.ID
}

// NewWorkerPool creates a new worker pool with x max workers and two channels
// for communication between the workers and the pool
func NewWorkerPool(max int) WorkerPool {
	return WorkerPool{
		MaxWorkers:        max,
		ProcessingChannel: make(chan bool, max),
		FinishedChannel:   make(chan worker.ID, max),
		ShutdownChannel:   make(chan bool, 1),
	}
}

// Start starts the worker pool and handles the communication
func (pool *WorkerPool) Start(handler func(worker worker.Worker)) {
	// infinite loop (until canceled) for communication
	for {
		select {
		case finished := <- pool.FinishedChannel:

		}
	}
}
