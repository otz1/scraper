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

// GoogleScraperImpl is an implemenation of a scraper
// service that will scrape google
type GoogleScraperImpl struct{}

// NewGoogleScraperService ...
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

func getSearchResultSet(query string) []*colly.HTMLElement {
	c := colly.NewCollector()
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
	return searchElements
}

func convertResults(searchElements []*colly.HTMLElement) []ScrapedResult {
	scrapedResults := []ScrapedResult{}

	for _, result := range searchElements {
		// todo get all the other stuff out of the result

		result.ForEach("a[href]", func(a int, e *colly.HTMLElement) {
			link := e.Attr("href")
			if !strings.HasPrefix(link, "/url?q=http") {
				return
			}
			text := e.ChildText("div")
			fmt.Println(e, text)

			result := ScrapedResult{}
			// TODO populate the scraped result.
			scrapedResults = append(scrapedResults, result)
		})
	}

	return scrapedResults
}

// Scrape will scrape google for the given query and
// parse the results.
func (g *GoogleScraperImpl) Scrape(query string) []ScrapedResult {
	log.Println("Scraping google for", query)
	resultSet := getSearchResultSet(query)
	convertedResults := convertResults(resultSet)
	return convertedResults
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
