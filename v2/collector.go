package main

import (
	"sync"
)

type Collector struct {
	UserAgent  string
	MaxDepth   string
	store      storage.Storage
	worker    *httpWorker
	wg        *sync.WaitGroup
	lock      *sync.RWMutex
}

func NewCollector() *Collector {
	c := &Collector{}
	c.Init()
}

func (c *Collector) Init() {
	c.UserAgent = "growler - https://github.com/Techassi/growler"
	c.MaxDepth = 0
	c.store = &storage.InMemory{}
	c.store.Init()
	c.worker = &httpWorker{}
	c.worker.Init()
	c.wg = &sync.WaitGroup{}
	c.lock = &sync.RWMutex{}
}

func (c *Collector) Visit(URL string) error {
	return c.build(URL, nil, false)
}

func (c *Collector) build(u string, depth int, revisit bool) error {
	c.checkRequest(u, depth, revisit)

	pURL, err := url.Parse(u)
	if err != nil {
		return err
	}

	c.wg.Add(1)
	go c.fetch(u, depth)
}

func (c *Collector) fetch(u string, depth int) error {
	defer c.wg.Done()

	c.worker.Request(u, depth)
}

func (c *Collector) checkRequest(u string, depth int, revisit bool) error {
	// Check if URL is empty. Throw ErrURLEmpty if so
	if u == "" {
		return ErrURLEmpty
	}

	// Check if depth is valid. Throw ErrDepthInvalid if not
	if (depth > 0 && depth > c.MaxDepth) || depth < 0 {
		return ErrDepthInvalid
	}

	// If we don't want to revisit the URL check if we already did. If so throw
	// ErrAlreadyVisited
	if !revisit {
		return ErrAlreadyVisited
	}

	return nil
}
