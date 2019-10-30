package crawl

import (
	"log"
	"strings"
	"net/http"

	"github.com/PuerkitoBio/goquery"

	"github.com/Techassi/growler/internal/queue"
)

func Crawl(data interface{}) (interface{}) {
	d := data.(queue.Job)
	res, err := http.Get(d.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
  	if err != nil {
    	log.Fatal(err)
  	}

	var links []string

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists && !strings.Contains(href, "mailto") {
			links = append(links, href)
		}
  	})

	return links
}
