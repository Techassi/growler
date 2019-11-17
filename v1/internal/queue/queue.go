package queue

import (
	"fmt"
	"errors"

	m "github.com/Techassi/growler/internal/models"
)

// NewQueue creates a new queue with N max items in it
func NewQueue(max int) (m.Queue, error) {
	if max < 0 {
		return m.Queue{}, errors.New(fmt.Sprintf("Provide a value greater than 0 for parameter 'max'"))
	}

	return m.Queue{
		MaxItems: max,
		Items: make([]m.Job, 0),
	}, nil
}
