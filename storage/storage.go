package storage

import (
	"net/url"
	"sync"
)

type InMemory struct {
	VisitedLinks map[*url.URL]bool
	lock         *sync.RWMutex
}

func (im *InMemory) Init() {
	im.VisitedLinks = make(map[*url.URL]bool)
	im.lock = &sync.RWMutex{}
}

func (im *InMemory) Visited(u *url.URL) {
	im.lock.Lock()
	im.VisitedLinks[u] = true
	im.lock.Unlock()
}

func (im *InMemory) IsVisited(u *url.URL) bool {
	im.lock.RLock()
	v := im.VisitedLinks[u]
	im.lock.RUnlock()

	return v
}
