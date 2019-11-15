package crawl

import (
	"log"
	"time"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	valid "github.com/asaskevich/govalidator"

	m "github.com/Techassi/growler/internal/models"
)

func Crawl(data interface{}, mode string) (interface{}) {
	d := data.(m.Job)

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

	if mode == "polite" {
		time.Sleep(2 * time.Second)
	}

	return links
}

func validURL(u string) (bool) {
	return valid.IsURL(u)
}
