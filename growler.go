package growler

import (
	"bytes"
	"errors"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/Techassi/growler/response"
	"github.com/Techassi/growler/storage"
)

type Collector struct {
	UserAgent       string
	MaxDepth        int
	Delay           int
	Duration        int
	startTime       time.Time
	store           *storage.InMemory
	worker          *httpWorker
	wg              *sync.WaitGroup
	lock            *sync.RWMutex
	onHTMLFunctions []Callback
}

type Callback struct {
	Function func(CollectorHTMLNode)
	Selector string
}

var (
	ErrURLEmpty         = errors.New("URL is empty")
	ErrDepthInvalid     = errors.New("Max depth limit reached or depth < 0")
	ErrAlreadyVisited   = errors.New("URL already visited")
	ErrHTTPStatusCode   = errors.New("HTTP status code of response is greater than 202")
	ErrDoubleSelector   = errors.New("A function with this selector was already registered")
	ErrReadingFromBody  = errors.New("Goquery couldn't read from body")
	ErrEmptyResponse    = errors.New("The body of the response is empty")
	ErrDepthExceeded    = errors.New("The depth of the current URL exceeds MaxDepth")
	ErrDurationExceeded = errors.New("The Duration is exceeded")
)

func NewCollector() *Collector {
	c := &Collector{}
	c.Init()

	return c
}

func (c *Collector) Init() {
	c.UserAgent = "growler - https://github.com/Techassi/growler"
	c.MaxDepth = 0
	c.Delay = 0
	c.Duration = 0
	c.store = &storage.InMemory{}
	c.store.Init()
	c.worker = &httpWorker{}
	c.worker.Init()
	c.wg = &sync.WaitGroup{}
	c.lock = &sync.RWMutex{}
}

func (c *Collector) Visit(URL string) error {
	if c.startTime.IsZero() {
		c.startTime = time.Now()
	}

	return c.build(URL, false)
}

func (c *Collector) Seeds(URLs []string) {
	for _, URL := range URLs {
		c.Visit(URL)
	}
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

func (c *Collector) SetDelay(d int) {
	if d <= 0 {
		return
	}

	c.Delay = d
}

func (c *Collector) SetMaxDepth(d int) {
	if d <= 0 {
		return
	}

	c.MaxDepth = d
}

func (c *Collector) SetDuration(d int) {
	if d <= 0 {
		return
	}

	c.Duration = d
}

func (c *Collector) Wait() {
	c.wg.Wait()
}

func (c *Collector) build(u string, revisit bool) error {
	pURL, err := c.checkRequest(u, revisit)
	if err != nil {
		return err
	}

	c.wg.Add(1)
	go c.fetch(pURL)

	return nil
}

func (c *Collector) fetch(u *url.URL) error {
	defer c.wg.Done()

	if c.Delay > 0 {
		time.Sleep(time.Duration(c.Delay) * time.Second)
	}

	res, err := c.worker.Request(u)
	if err != nil {
		return err
	}

	err = c.checkHTTPStatusCode(res)
	if err != nil {
		return err
	}

	err = c.handleHTML(res)

	c.store.Visited(u)

	return nil
}

func (c *Collector) checkRequest(u string, revisit bool) (*url.URL, error) {
	if time.Now().Sub(c.startTime).Seconds() > float64(c.Duration) && c.Duration > 0 {
		return nil, ErrDurationExceeded
	}

	if u == "" {
		return nil, ErrURLEmpty
	}

	if c.MaxDepth < 0 {
		return nil, ErrDepthInvalid
	}

	pURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	if !revisit && c.store.IsVisited(pURL) {
		return nil, ErrAlreadyVisited
	}

	if strings.Count(pURL.Path, "/") > c.MaxDepth && c.MaxDepth > 0 {
		return nil, ErrDepthExceeded
	}

	return pURL, nil
}

func (c *Collector) checkHTTPStatusCode(res *response.Response) error {
	if res == nil {
		return ErrEmptyResponse
	}

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
					Name:       n.Data,
					Collector:  c,
					attributes: n.Attr,
				})
			}
		})
	}

	return nil
}
