package models

type Event struct {
	Type    string
	Worker  Worker
	Pool   *WorkerPool
}
