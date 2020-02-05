package growler

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/Techassi/growler/response"
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

func (h *httpWorker) Request(u *url.URL) (*response.Response, error) {
	r := &http.Request{
		Method:     "GET",
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	res, err := h.Client.Do(r)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	b, e := ioutil.ReadAll(res.Body)
	if e != nil {
		return nil, err
	}

	return &response.Response{
		StatusCode: res.StatusCode,
		Body:       b,
		Headers:    res.Header,
	}, nil
}
