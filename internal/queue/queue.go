package queue

import (
	"errors"
)

type Queue struct {
	Items []string
}

func (queue Queue) Poll() (string, error) {
	// check if there are items in the queue
	if len(queue.Items == 0) {
		return "", errors.New("Queue is empty")
	}

	// get item to process
	item = queue.Items[0]

	// truncate the array
	queue.Items[1:]

	return item, nil
}
