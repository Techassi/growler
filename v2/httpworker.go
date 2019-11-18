package main

import (
	"sync"
	"io/ioutil"
	"net/url"
	"net/http"

	"github.com/Techassi/growler/v2/response"
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
	}

	res, err := h.CLient.Do(r)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return response.Response{
		StatusCode: res.StatusCode,
		Body: ioutil.ReadAll(res.Body),
		Header: res.Header,
	}, nil
}
