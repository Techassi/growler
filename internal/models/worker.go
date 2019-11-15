package models

import (
	"github.com/google/uuid"
)

type Worker struct {
	ID                uuid.UUID
	Mode              string
	Action 	          func(interface{}, string) interface{}
	JobChannel        <-chan Job
	ResultChannel     chan<- interface{}
	LifecycleChannel  chan<- Event
}

func NewWorker(jC <-chan Job, rC chan<- interface{}, lC chan<- Event, ac func(interface{}, string) interface{}, mode string) Worker {
	return Worker{
		ID:                uuid.New(),
		Mode:              mode,
		Action:            ac,
		JobChannel:        jC,
		ResultChannel:     rC,
		LifecycleChannel:  lC,
	}
}

func (worker Worker) Run() {
	// Lifecycle worker:init
	worker.LifecycleChannel <- Event{
		Type: "worker:init",
		Worker: worker,
	}

	for job := range worker.JobChannel {
		// Lifecycle worker:processing
		worker.LifecycleChannel <- Event{
			Type: "worker:process",
			Worker: worker,
		}

		// do the actual work
		worker.ResultChannel <- worker.Action(job, worker.Mode)

		// Lifecycle worker:finished
		worker.LifecycleChannel <- Event{
			Type: "worker:finish",
			Worker: worker,
		}
	}
}
