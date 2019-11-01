package queue

import (
	"fmt"
	"errors"
)

type Queue struct {
	MaxItems   int
	Items    []Job
}

type Job struct {
	Priority int
	URL      string
}

// NewQueue creates a new queue with N max items in it
func NewQueue(max int) (Queue, error) {
	if max < 0 {
		return Queue{}, errors.New(fmt.Sprintf("Provide a value greater than 0 for parameter 'max'"))
	}

	return Queue{
		MaxItems: max,
		Items: make([]Job, 0),
	}, nil
}

// Poll takes the first element in queue and returns it or returns an error if
// the queue is empty
// TODO: This is 9/10 a huge performance bottleneck
func (queue *Queue) Poll() (Job, error) {
	// check if there are items in the queue
	if len(queue.Items) == 0 {
		return Job{}, errors.New("Queue is empty")
	}

	// get item to process
	item := queue.Items[0]

	// truncate the array
	queue.Items = queue.Items[1:]

	return item, nil
}

// Queue queues a new job in the queue. The new job gets appended at the end of
// the queue
func (queue *Queue) Queue(job Job) (Job, error) {
	if len(queue.Items) > queue.MaxItems {
		return job, errors.New(fmt.Sprintf("Queue is full (MaxItems: %d). Returning job for later queueing", queue.MaxItems))
	}

	queue.Items = append(queue.Items, job)
	return job, nil
}

func (queue *Queue) URLJob(url_string string) error {
	_, err := queue.Queue(Job{
		Priority: 1,
		URL: url_string,
	})
	if err != nil {
		return err
	}

	return nil
}

func (queue *Queue) QueueList(l interface{}) {
	links := l.([]string)

	for _, link := range links {
		queue.Queue(Job{
			Priority: 1,
			URL: link,
		})
	}
}
