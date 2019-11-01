package workerpool

import (
	"fmt"
	"time"
	"errors"

	"github.com/google/uuid"

	"github.com/Techassi/growler/internal/queue"
	"github.com/Techassi/growler/internal/worker"
	"github.com/Techassi/growler/internal/helper"
)

type WorkerPool struct {
	Queue             queue.Queue
	Action            func(interface{}) interface{}
	MaxWorkers        int
	Events            map[string] func(*WorkerPool)
	ActiveWorkers     map[uuid.UUID]time.Time
	ShutdownChannel   chan bool
	JobChannel        chan queue.Job
	ResultChannel     chan interface{}
	LifecycleChannel  chan string
}

type Event struct {
	Type       string
	Worker    *worker.Worker
	Pool      *WorkerPool
}

type Config struct {
	Verbose bool
}

// NewWorkerPool creates a new worker pool with N max workers and two channels
// for communication between the workers and the pool
func NewWorkerPool(max int, q queue.Queue, action func(interface{}) interface{}) (WorkerPool, error) {
	if max < 0 {
		return WorkerPool{}, errors.New(fmt.Sprintf("Provide a value greater than 0 for parameter 'max'"))
	}

	return WorkerPool{
		Queue: q,
		Action: action,
		MaxWorkers:        max,
		Events:            make(map[string] func(*WorkerPool)),
		ActiveWorkers:     make(map[uuid.UUID]time.Time, max),
		ShutdownChannel:   make(chan bool, 1),
		JobChannel:        make(chan queue.Job, max),
		ResultChannel:     make(chan interface{}, max),
		LifecycleChannel:  make(chan string, max),
	}, nil
}

// Start starts the worker pool and handles the communication
func (pool *WorkerPool) Start() {
	// setup workers based on pool.MaxWorkers
	for i := 0; i < pool.MaxWorkers; i++ {
		worker := worker.NewWorker(pool.JobChannel, pool.ResultChannel, pool.LifecycleChannel, pool.Action)
		pool.AddWorker(worker.ID)

		go worker.Run()
	}

	// infinite loop (until canceled) for communication
	for {
		select {
		case result := <-pool.ResultChannel:
			pool.Queue.QueueList(result)
		case event := <-pool.LifecycleChannel:
			pool.do(event)
		case shutdown := <-pool.ShutdownChannel:
			fmt.Printf("Shutdown %t. Exiting...", shutdown)
			break
		}
	}
}

// AddWorker adds a new worker and ActiveWorkers keeps track of all active
// workers
func (pool *WorkerPool) AddWorker(id uuid.UUID) (error) {
	if _, ok := pool.ActiveWorkers[id]; ok {
		pool.ActiveWorkers[id] = time.Now()

		return nil
	}

	return errors.New(fmt.Sprintf("There is already a worker with id %v", id))
}

// On registers an event function getting triggered when x event is called by
// LifecycleChannel
func (pool *WorkerPool) On(event string, action func(*WorkerPool)) (error) {
	if e, ok := pool.Events[event]; ok {
		return errors.New(fmt.Sprintf("Already added function for event %s: %s", event, helper.GetFunctionName(e)))
	}

	pool.Events[event] = action
	return nil
}

// Executes a provided function. Used to trigger event function on
// Lifecycle events
func (pool *WorkerPool) do(event string) {
	e, ok := pool.Events[event]

	// default event functions
	switch event {
	case "worker:init":
		// fmt.Println("worker:init")

		// When one worker is initialized poll the queue for a new job
		// and push it into the JobChannel
		poll, err := pool.Queue.Poll()
		if err == nil {
			pool.JobChannel <- poll
		}
	case "worker:processing":
		// fmt.Println("worker:processing")
	case "worker:finished":
		// fmt.Println("worker:finished")

		// When one worker finished his work poll the queue for a new job
		// and push it into the JobChannel
		poll, err := pool.Queue.Poll()
		if err == nil {
			pool.JobChannel <- poll
		}
	}

	// execute custom event function registered for event
	if ok { e(pool) }
}
