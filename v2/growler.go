package growler

import (
	"sync"
	"bytes"
	"errors"
	"net/url"

	"github.com/PuerkitoBio/goquery"

	"github.com/Techassi/growler/v2/storage"
	"github.com/Techassi/growler/v2/response"
)

type Collector struct {
	UserAgent         string
	MaxDepth          int
	store            *storage.InMemory
	worker           *httpWorker
	wg               *sync.WaitGroup
	lock             *sync.RWMutex
	onHTMLFunctions []Callback
}

type Callback struct {
	Function func (CollectorHTMLNode)
	Selector string
}

var (
	ErrURLEmpty        = errors.New("URL is empty")
	ErrDepthInvalid    = errors.New("Max depth limit reached or depth < 0")
	ErrAlreadyVisited  = errors.New("URL already visited")
	ErrHTTPStatusCode  = errors.New("HTTP status code of response is greater than 202")
	ErrDoubleSelector  = errors.New("A function with this selector was already registered")
	ErrReadingFromBody = errors.New("Goquery couldn't read from body")
)

func NewCollector() *Collector {
	c := &Collector{}
	c.Init()

	return c
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
	return c.build(URL, 0, false)
}

func (c *Collector) OnHTML(selector string, f func(CollectorHTMLNode)) {
	c.lock.Lock()

	if c.onHTMLFunctions == nil {
		c.onHTMLFunctions = make([]Callback, 0, 5)
	}

	c.onHTMLFunctions = append(c.onHTMLFunctions, Callback{
		Function: f,
		Selector: selector,
	})

	c.lock.Unlock()
}

func (c *Collector) Wait() {
	c.wg.Wait()
}

func (c *Collector) build(u string, depth int, revisit bool) error {
	c.checkRequest(u, depth, revisit)

	pURL, err := url.Parse(u)
	if err != nil {
		return err
	}

	c.wg.Add(1)
	go c.fetch(pURL, depth)

	return nil
}

func (c *Collector) fetch(u *url.URL, depth int) error {
	defer c.wg.Done()

	res, err := c.worker.Request(u, depth)
	if err != nil {
		return err
	}

	err = c.checkHTTPStatusCode(res)
	if err != nil {
		return err
	}

	err = c.handleHTML(res)

	return nil
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
	if !revisit && c.store.IsVisited(u) {
		return ErrAlreadyVisited
	}

	return nil
}

func (c *Collector) checkHTTPStatusCode(res *response.Response) error {
	if res.StatusCode < 203 {
		return nil
	}

	return ErrHTTPStatusCode
}

func (c *Collector) handleHTML(res *response.Response) error {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(res.Body))
	if err != nil {
		return err
	}

	for _, call := range c.onHTMLFunctions {
		doc.Find(call.Selector).Each(func(_ int, s *goquery.Selection) {
			for _, n := range s.Nodes {
				call.Function(CollectorHTMLNode{
					Name: n.Data,
					Collector: c,
					attributes: n.Attr,
				})
			}
	  	})
	}

	return nil
}
