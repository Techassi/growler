package worker

import (
	"fmt"
	
	"github.com/google/uuid"

	"github.com/Techassi/growler/internal/queue"
)

type Worker struct {
	ID                uuid.UUID
	ProcessingChannel chan<- bool
	FinishedChannel   chan<- uuid.UUID
	JobChannel        <-chan queue.Job
}

func NewWorker(pC chan<- bool, fC chan<- uuid.UUID, jC <-chan queue.Job) Worker {
	return Worker{
		ID:                uuid.New(),
		ProcessingChannel: pC,
		FinishedChannel:   fC,
		JobChannel:        jC,
	}
}

func (worker Worker) Run() {
	worker.ProcessingChannel <- true

	for job := range worker.JobChannel {
		worker.RunFunc(job)
		worker.FinishedChannel <- worker.ID
	}
}

func (worker Worker) RunFunc(job queue.Job) {
	fmt.Println(job.URL)
}
