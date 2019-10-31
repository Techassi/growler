package workerpool

import (
	"os"
	"fmt"
	"time"
	"errors"
	"reflect"
	"runtime"

	"github.com/google/uuid"

	"github.com/Techassi/growler/internal/worker"
	"github.com/Techassi/growler/internal/queue"
)

var counted int
var indexed []string

type WorkerPool struct {
	Queue             queue.Queue
	Action            func(interface{}) interface{}
	MaxWorkers        int
	Events            map[string] func(*WorkerPool)
	ActiveWorkers     map[uuid.UUID]time.Time
	ProcessingChannel chan bool
	FinishedChannel   chan uuid.UUID
	ShutdownChannel   chan bool
	JobChannel        chan queue.Job
	ResultChannel     chan interface{}
	LifecycleChannel  chan string
}

// NewWorkerPool creates a new worker pool with x max workers and two channels
// for communication between the workers and the pool
func NewWorkerPool(max int, q queue.Queue, action func(interface{}) interface{}) WorkerPool {
	return WorkerPool{
		Queue: q,
		Action: action,
		MaxWorkers:        max,
		Events:            make(map[string] func(*WorkerPool)),
		ActiveWorkers:     make(map[uuid.UUID]time.Time, max),
		ProcessingChannel: make(chan bool, max),
		FinishedChannel:   make(chan uuid.UUID, max),
		ShutdownChannel:   make(chan bool, 1),
		JobChannel:        make(chan queue.Job, max),
		ResultChannel:     make(chan interface{}, max),
		LifecycleChannel:  make(chan string, max),
	}
}

// Start starts the worker pool and handles the communication
func (pool *WorkerPool) Start() {
	// setup workers based on pool.MaxWorkers
	for i := 0; i < pool.MaxWorkers; i++ {
		worker := worker.NewWorker(pool.ProcessingChannel, pool.FinishedChannel, pool.JobChannel, pool.ResultChannel, pool.LifecycleChannel, pool.Action)
		pool.AddWorker(worker.ID)

		go worker.Run()
	}

	// infinite loop (until canceled) for communication
	for {
		select {
		case finished := <-pool.FinishedChannel:
			_ = finished
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
			counted += len(result.([]string))
			for _, r := range result.([]string) {
				indexed = append(indexed, r)
			}

			if counted > 10000 {
				f, err := os.Create("data.txt")
				if err != nil {
					panic(err)
				}
				defer f.Close()

				for _, i := range indexed {
					w, err := f.WriteString(fmt.Sprintf("%s\n", i))
					if err != nil {
						panic(err)
						fmt.Printf("wrote %d bytes\n", w)
					}
				}
				f.Sync()
				os.Exit(1)
			}
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
		return errors.New(fmt.Sprintf("Already added function for event %s: %s", event, getFunctionName(e)))
	}

	pool.Events[event] = action
	return nil
}

// Executes a provided function. Used to trigger event function on
// Lifecycle events
func (pool *WorkerPool) do(event string) {
	if e, ok := pool.Events[event]; ok {
		e(pool)
	}
}

func getFunctionName(i interface{}) string {
    return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
