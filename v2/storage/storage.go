package storage

import (
	"time"
	"sync"
)

type InMemory struct {
	Visited  map[string]bool
	lock    *sync.RWMutex
}

func (im *InMemory) Init() {
	im.Visited = make(map[string]bool)
	im.lock = &sync.RWMutex{}
}

func (im *InMemory) Visited(u string) {
	im.lock.Lock()
	s.Visited[u] = true
	im.lock.Unlock()
}

func (im *InMemory) IsVisited(u string) {
	s.lock.RLock()
	v := im.Visited[u]
	s.lock.RUnlock()

	return v
}
