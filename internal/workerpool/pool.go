package workerpool

import (
	// "os"
	"fmt"
	"time"
	"errors"

	"github.com/google/uuid"

	"github.com/Techassi/growler/internal/worker"
	"github.com/Techassi/growler/internal/queue"
)

type WorkerPool struct {
	Queue             queue.Queue
	Action            func(interface{}) interface{}
	MaxWorkers        int
	ActiveWorkers     map[uuid.UUID]time.Time
	ProcessingChannel chan bool
	FinishedChannel   chan uuid.UUID
	ShutdownChannel   chan bool
	JobChannel        chan queue.Job
	ResultChannel     chan interface{}
}

// NewWorkerPool creates a new worker pool with x max workers and two channels
// for communication between the workers and the pool
func NewWorkerPool(max int, q queue.Queue, action func(interface{}) interface{}) WorkerPool {
	return WorkerPool{
		Queue: q,
		Action: action,
		MaxWorkers:        max,
		ActiveWorkers:     make(map[uuid.UUID]time.Time, max),
		ProcessingChannel: make(chan bool, max),
		FinishedChannel:   make(chan uuid.UUID, max),
		ShutdownChannel:   make(chan bool, 1),
		JobChannel:        make(chan queue.Job, max),
		ResultChannel:     make(chan interface{}, max),
	}
}

// Start starts the worker pool and handles the communication
func (pool *WorkerPool) Start() {
	// setup workers based on pool.MaxWorkers
	for i := 0; i < pool.MaxWorkers; i++ {
		worker := worker.NewWorker(pool.ProcessingChannel, pool.FinishedChannel, pool.JobChannel, pool.ResultChannel, pool.Action)
		pool.AddWorker(worker.ID)

		go worker.Run()
	}

	// poll queue to start
	// poll, err := pool.Queue.Poll()
	// if err == nil {
	// 	pool.JobChannel <- poll
	// }

	// infinite loop (until canceled) for communication
	for {
		select {
		case finished := <-pool.FinishedChannel:
			fmt.Printf("Finished | Worker %v\n", finished)
			// When one worker finished his work poll the queue for a new job
			// and push it into the JobChannel
			poll, err := pool.Queue.Poll()
			if err == nil {
				pool.JobChannel <- poll
			}
		case processing := <-pool.ProcessingChannel:
			if !processing {
				fmt.Println(processing)
			}
		case result := <-pool.ResultChannel:
			pool.Queue.QueueList(result)
		case shutdown := <-pool.ShutdownChannel:
			fmt.Printf("Shutdown %t. Exiting...", shutdown)
			break
		}
	}
}

func (pool *WorkerPool) AddWorker(id uuid.UUID) (error) {
	if _, ok := pool.ActiveWorkers[id]; ok {
		pool.ActiveWorkers[id] = time.Now()

		return nil
	}

	return errors.New(fmt.Sprintf("There is already a worker with id %v", id))
}
