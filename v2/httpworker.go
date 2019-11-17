package main

import (
	"sync"
	"net/url"
	"net/http"
)

type httpWorker struct {
	Client *http.Client
	lock   *sync.RWMutex
}

func (h *httpWorker) Init() {
	h.Client = &http.Client{
		Timeout: 10 * time.Second,
	}
	h.lock = &sync.RWMutex{}
}

func (h *httpWorker) Request(u url.URL, depth int) (*Response, error) {
	r := http.Request{
		Method: "GET",
		URL: u,
		Proto: "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		// Body and Response missing
	}

	h.CLient.Do(r)
}
