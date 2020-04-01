package scraper

// https://edmundmartin.com/scraping-google-with-golang/

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

// GoogleScraperImpl is an impl. of a scraper
// service that will scrape google
type GoogleScraperImpl struct{}

// NewGoogleScraperService ...
func NewGoogleScraperService() *GoogleScraperImpl {
	return &GoogleScraperImpl{}
}

func (g *GoogleScraperImpl) buildRequestURL(query, langCode string) string {
	tld := map[string]string{
		"us": "com",
		"gb": "co.uk",
		"ru": "ru",
		"fr": "fr",
	}[langCode]

	return fmt.Sprintf(`https://www.google.%s/search?q=%s&hl=%s`, tld, parseQuery(query), langCode)
}

func (g *GoogleScraperImpl) getSearchResultSet(query string) []*colly.HTMLElement {
	log.Println("getSearchResultSet for", query)

	c := colly.NewCollector()
	var searchElements []*colly.HTMLElement
	c.OnHTML("div", func(e *colly.HTMLElement) {

		e.ForEach("#search", func(count int, e *colly.HTMLElement) {
			log.Println(e)
		})

		// searchElements = append(searchElements, e)
	})
	url := g.buildRequestURL(query, "gb")
	if err := c.Visit(url); err != nil {
		panic(err)
	}
	return searchElements
}

func (g *GoogleScraperImpl) convertResults(searchElements []*colly.HTMLElement) []ScrapedResult {
	var scrapedResults []ScrapedResult

	for _, result := range searchElements {
		fmt.Println("result is", result.DOM.Closest("a[href]").Text())

		result.ForEach("a[href]", func(a int, e *colly.HTMLElement) {
			link := e.Attr("href")

			// its not a url that we care about
			if !strings.HasPrefix(link, "/url?q=http") {
				return
			}

			// trim the prefix
			link = strings.TrimPrefix(link, "/url?q=")

			titles := e.ChildTexts("h3")
			title := "no title!"
			if len(titles) > 0 {
				title = titles[0]
			}

			// ensure it's a proper link!
			if _, err := url.Parse(link); err != nil {
				log.Println("failed to parse url", err)
				return
			}

			log.Println(title, "->", link)

			result := ScrapedResult{
				Href: link,
			}
			scrapedResults = append(scrapedResults, result)
		})
	}

	return scrapedResults
}

// Scrape will scrape google for the given query and
// parse the results.
func (g *GoogleScraperImpl) Scrape(query string) []ScrapedResult {
	log.Println("Scraping google for", query)
	resultSet := g.getSearchResultSet(query)
	convertedResults := g.convertResults(resultSet)
	return convertedResults
}

func (g *GoogleScraperImpl) getPageContents(query, langCode string) string {
	url := g.buildRequestURL(query, langCode)
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
