package scraper

// https://edmundmartin.com/scraping-google-with-golang/

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
)

type GoogleScraperImpl struct {
}

func NewGoogleScraperService() *GoogleScraperImpl {
	return &GoogleScraperImpl{}
}

func buildRequestURL(query, langCode string) string {
	tld := map[string]string{
		"us": "com",
		"gb": "co.uk",
		"ru": "ru",
		"fr": "fr",
	}[langCode]

	return fmt.Sprintf(`https://www.google.%s/search?q=%s&hl=%s`, tld, parseQuery(query), langCode)
}

func (g *GoogleScraperImpl) Scrape(query string) []ScrapedResult {
	c := colly.NewCollector()
	log.Println("Scraping google for", query)

	searchElements := []*colly.HTMLElement{}
	c.OnHTML("div", func(e *colly.HTMLElement) {
		e.ForEach("div", func(a int, el *colly.HTMLElement) {
			classList := strings.Split(el.Attr("class"), " ")
			if len(classList) == 4 {
				searchElements = append(searchElements, el)
			}
		})
	})
	url := buildRequestURL(query, "gb")
	c.Visit(url)

	scrapedResults := []ScrapedResult{}

	for _, result := range searchElements {
		result.ForEach("a[href]", func(a int, e *colly.HTMLElement) {
			link := e.Attr("href")
			if !strings.HasPrefix(link, "/url?q=http") {
				return
			}
			text := e.ChildText("div")

			scrapedResults = append(scrapedResults, ScrapedResult{
				Title:   text,
				Href:    link,
				Keyword: []string{},
			})
		})
	}

	return scrapedResults
}

func getPageContents(query, langCode string) string {
	url := buildRequestURL(query, langCode)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(data)
}
