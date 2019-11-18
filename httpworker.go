package growler

import (
	"sync"
	"time"
	"io/ioutil"
	"net/url"
	"net/http"

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

func (h *httpWorker) Request(u *url.URL, delay int) (*response.Response, error) {
	if delay > 0 {
		time.Sleep(time.Duration(delay) * time.Second)
	}

	r := &http.Request{
		Method: "GET",
		URL: u,
		Proto: "HTTP/1.1",
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
		Body: b,
		Headers: res.Header,
	}, nil
}
