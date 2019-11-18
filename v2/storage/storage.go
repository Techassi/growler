package storage

import (
	"sync"
)

type InMemory struct {
	VisitedLinks  map[string]bool
	lock         *sync.RWMutex
}

func (im *InMemory) Init() {
	im.VisitedLinks = make(map[string]bool)
	im.lock = &sync.RWMutex{}
}

func (im *InMemory) Visited(u string) {
	im.lock.Lock()
	im.VisitedLinks[u] = true
	im.lock.Unlock()
}

func (im *InMemory) IsVisited(u string) bool {
	im.lock.RLock()
	v := im.VisitedLinks[u]
	im.lock.RUnlock()

	return v
}
