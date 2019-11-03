package crawl

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	valid "github.com/asaskevich/govalidator"

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
		if href, exists := s.Attr("href"); exists && validURL(href) {
			links = append(links, href)
		}
  	})

	return links
}

func validURL(u string) (bool) {
	return valid.IsURL(u)
}
