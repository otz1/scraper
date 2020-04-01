package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// DuckDuckGoScraperImpl is an implementation a scraper
// service that will scrape duckduckgo
type DuckDuckGoScraperImpl struct{}

// NewDDGScraperService ...
func NewDDGScraperService() *DuckDuckGoScraperImpl {
	return &DuckDuckGoScraperImpl{}
}

func (d *DuckDuckGoScraperImpl) buildRequestURL(query, langCode string) string {
	tld := map[string]string{
		"us": "com",
		"gb": "co.uk",
	}[langCode]

	return fmt.Sprintf(`https://duckduckgo.%s/html/?q=%s`, tld, parseQuery(query))
}

func (d *DuckDuckGoScraperImpl) convertLink(link string) (string, bool) {
	decodedLink, err := url.QueryUnescape(link)
	if err != nil {
		log.Fatal(err)
		return "", false
	}
	idx := strings.Index(decodedLink, "http")
	if idx == -1 {
		return "", false
	}
	return decodedLink[idx:], true
}

func (d *DuckDuckGoScraperImpl) getSearchResultSet(query string) []ScrapedResult {
	log.Println("getSearchResultSet for", query)

	c := colly.NewCollector()
	var results []ScrapedResult

	// here we are scraping the container of all of the result boxes.
	c.OnHTML("div[id=links]", func(e *colly.HTMLElement) {
		e.DOM.Find(".result").Each(func(_ int, s *goquery.Selection) {
			body := s.Find(".result__body")
			title := body.Find(".result__title a").Text()

			link, ok := body.Find(".result__title .result__a").Attr("href")
			if !ok {
				log.Println("failed to get URL for result", title)
				return
			}

			convertedLink, ok := d.convertLink(link)
			if !ok {
				log.Println("failed to convert URL for result", title)
				return
			}

			log.Println(convertedLink, "->", title)

			results = append(results, ScrapedResult{
				Title: title,
				Href: convertedLink,
			})
		})
	})

	url := d.buildRequestURL(query, "gb")
	if err := c.Visit(url); err != nil {
		panic(err)
	}

	return results
}

// Scrape will scrape google for the given query and
// parse the results.
func (d *DuckDuckGoScraperImpl) Scrape(query string) []ScrapedResult {
	log.Println("Scraping ddg for", query)
	convertedResults := d.getSearchResultSet(query)
	return convertedResults
}

func (d *DuckDuckGoScraperImpl) getPageContents(query, langCode string) string {
	url := d.buildRequestURL(query, langCode)
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
