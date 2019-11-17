package workerpool

import (
	"fmt"
	"time"
	"errors"

	"github.com/google/uuid"

	m "github.com/Techassi/growler/internal/models"
)

// NewWorkerPool creates a new worker pool with N max workers and four channels
// for communication between the workers and the pool
func NewWorkerPool(max int, q m.Queue, action func(interface{}, string) interface{}) (m.WorkerPool, error) {
	if max < 0 {
		return m.WorkerPool{}, errors.New(fmt.Sprintf("Provide a value greater than 0 for parameter 'max'"))
	}

	return m.WorkerPool{
		Queue: q,
		Action: action,
		MaxWorkers:        max,
		Mode:			   "polite",
		Events:            make(map[string] func(m.Event)),
		ActiveWorkers:     make(map[uuid.UUID]time.Time, max),
		ShutdownChannel:   make(chan bool, 1),
		JobChannel:        make(chan m.Job, max),
		ResultChannel:     make(chan interface{}, max),
		LifecycleChannel:  make(chan m.Event, max),
	}, nil
}
