package worker

import (
	"github.com/google/uuid"
)

type Worker struct {
	ID                uuid.UUID
	ProcessingChannel chan<- bool
	FinishedChannel   chan<- uuid.UUID
}

func NewWorker(processingChannel chan<- bool, finishedChannel chan<- ID) Worker {
	return Worker{
		ID:                uuid.New(),
		ProcessingChannel: processingChannel,
		FinishedChannel:   finishedChannel,
	}
}
