package models

import (
	"fmt"
	"time"
	"errors"
	"github.com/google/uuid"

	"github.com/Techassi/growler/internal/helper"
)

type WorkerPool struct {
	Queue             Queue
	Action            func(interface{}, string) interface{}
	MaxWorkers        int
	Mode              string
	Events            map[string] func(Event)
	ActiveWorkers     map[uuid.UUID]time.Time
	ShutdownChannel   chan bool
	JobChannel        chan Job
	ResultChannel     chan interface{}
	LifecycleChannel  chan Event
}

type Config struct {
	Verbose bool
}

// Start starts the worker pool and handles the communication
func (pool *WorkerPool) Start() {
	// Lifecycle pool:init
	pool.LifecycleChannel <- Event{
		Type: "pool:init",
	}

	// setup workers based on pool.MaxWorkers
	for i := 0; i < pool.MaxWorkers; i++ {
		worker := NewWorker(pool.JobChannel, pool.ResultChannel, pool.LifecycleChannel, pool.Action, pool.Mode)
		pool.AddWorker(worker.ID)

		// Lifecycle pool:addworker
		// pool.LifecycleChannel <- "pool:addworker"

		go worker.Run()
	}

	// infinite loop (until canceled) for communication
	for {
		select {
		case result := <-pool.ResultChannel:
			pool.Queue.QueueList(result)
		case event := <-pool.LifecycleChannel:
			pool.do(Event{
				Type: event.Type,
				Worker: event.Worker,
				Pool: pool,
			})
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
func (pool *WorkerPool) On(event string, action func(Event)) (error) {
	if e, ok := pool.Events[event]; ok {
		return errors.New(fmt.Sprintf("Already added function for event %s: %s", event, helper.GetFunctionName(e)))
	}

	pool.Events[event] = action
	return nil
}

func (pool *WorkerPool) SetMode(mode string) (error) {
	switch mode {
	case "polite":
		pool.Mode = mode
		return nil
	case "speed":
		pool.Mode = mode
		return nil
	default:
		pool.Mode = "polite"
		return errors.New(fmt.Sprintf("%s is a not supported mode. Defaulted to 'polite'", mode))
	}
}

// Executes a provided function. Used to trigger event function on
// Lifecycle events
func (pool *WorkerPool) do(event Event) {
	e, ok := pool.Events[event.Type]

	// default event functions
	switch event.Type {
	case "worker:init":
		// When one worker is initialized poll the queue for a new job
		// and push it into the JobChannel
		poll, err := pool.Queue.Poll()
		if err == nil {
			pool.JobChannel <- poll
		}
	case "worker:process":
		// fmt.Println("worker:processing")
	case "worker:finish":
		// When one worker finished his work poll the queue for a new job
		// and push it into the JobChannel
		poll, err := pool.Queue.Poll()
		if err == nil {
			pool.JobChannel <- poll
		}
	case "pool:init":
		fmt.Println("pool:init")
	case "pool:addworker":
		fmt.Println("pool:addworker")
	}

	// execute custom event function registered for event
	if ok { e(event) }
}
