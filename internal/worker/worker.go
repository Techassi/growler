package worker

import (
	"github.com/google/uuid"

	"github.com/Techassi/growler/internal/queue"
)

type Worker struct {
	ID                uuid.UUID
	Action 	          func(interface{}) interface{}
	JobChannel        <-chan queue.Job
	ResultChannel     chan<- interface{}
	LifecycleChannel  chan<- string
}

func NewWorker(jC <-chan queue.Job, rC chan<- interface{}, lC chan<- string, ac func(interface{}) interface{}) Worker {
	return Worker{
		ID:                uuid.New(),
		Action:            ac,
		JobChannel:        jC,
		ResultChannel:     rC,
		LifecycleChannel:  lC,
	}
}

func (worker Worker) Run() {
	// Lifecycle worker:init
	worker.LifecycleChannel <- "worker:init"

	for job := range worker.JobChannel {
		// Lifecycle worker:processing
		worker.LifecycleChannel <- "worker:processing"

		// do work
		worker.ResultChannel <- worker.Action(job)

		// Lifecycle worker:finished
		worker.LifecycleChannel <- "worker:finished"
	}
}
