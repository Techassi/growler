package worker

import (
	"github.com/google/uuid"

	"github.com/Techassi/growler/internal/queue"
)

type Worker struct {
	ID                uuid.UUID
	Action 	          func(interface{}) interface{}
	ProcessingChannel chan<- bool
	FinishedChannel   chan<- uuid.UUID
	JobChannel        <-chan queue.Job
	ResultChannel     chan<- interface{}
}

func NewWorker(pC chan<- bool, fC chan<- uuid.UUID, jC <-chan queue.Job, rC chan<- interface{}, action func(interface{}) interface{}) Worker {
	return Worker{
		ID:                uuid.New(),
		Action:            action,
		ProcessingChannel: pC,
		FinishedChannel:   fC,
		JobChannel:        jC,
		ResultChannel:     rC,
	}
}

func (worker Worker) Run() {
	// declare as finished to initiate poll
	worker.FinishedChannel <- worker.ID

	for job := range worker.JobChannel {
		worker.ProcessingChannel <- true
		worker.ResultChannel <- worker.Action(job)
		worker.FinishedChannel <- worker.ID
	}
}
