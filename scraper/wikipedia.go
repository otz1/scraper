package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/getsentry/sentry-go"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// WikipediaScraperImpl is an implementation a scraper
// service that will scrape wikipedia
type WikipediaScraperImpl struct{}

// NewWikipediaScraperService ...
func NewWikipediaScraperService() *WikipediaScraperImpl {
	return &WikipediaScraperImpl{}
}

func (w *WikipediaScraperImpl) getBaseLink(subdomain string) string {
	return fmt.Sprintf("https://%s.wikipedia.org", subdomain)
}

func (w *WikipediaScraperImpl) buildRequestURL(query, langCode string) string {
	// TODO get a few other sources
	subdomain := map[string]string{
		"us": "en",
		"gb": "en",
	}[langCode]

	baseLink := w.getBaseLink(subdomain)
	return fmt.Sprintf(`%s/wiki/index.php?search=/html/?q=%s&profile=default&fulltext=1&ns0=1`, baseLink, parseQuery(query))
}

func (w *WikipediaScraperImpl) convertLink(link string) (string, bool) {
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

func (w *WikipediaScraperImpl) getSearchResultSet(query string) []ScrapedResult {
	log.Println("getSearchResultSet for", query)

	c := colly.NewCollector()
	var results []ScrapedResult

	// here we are scraping the container of all of the result boxes.
	c.OnHTML("div .mw-search-results", func(e *colly.HTMLElement) {
		e.DOM.Find(".mw-search-result").Each(func(_ int, body *goquery.Selection) {
			title := body.Find(".mw-search-result-heading a").Text()
			meta := body.Find(".searchresult").Text()

			pageLink, ok := body.Find(".mw-search-result-heading a").Attr("href")
			if !ok {
				log.Println("failed to get URL for result", title)
				return
			}

			link := w.articleLinkToAbsLink(pageLink)

			convertedLink, ok := w.convertLink(link)
			if !ok {
				log.Println("failed to convert URL for result", title)
				return
			}

			log.Println(convertedLink, "->", title)

			results = append(results, ScrapedResult{
				Title: title,
				Href: convertedLink,
				Snippet: meta,
			})
		})
	})

	url := w.buildRequestURL(query, "gb")
	if err := c.Visit(url); err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	return results
}

// Scrape will scrape google for the given query and
// parse the results.
func (w *WikipediaScraperImpl) Scrape(query string) []ScrapedResult {
	log.Println("Scraping wikipedia for", query)
	convertedResults := w.getSearchResultSet(query)
	return convertedResults
}

func (w *WikipediaScraperImpl) getPageContents(query, langCode string) string {
	url := w.buildRequestURL(query, langCode)
	resp, err := http.Get(url)
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	return string(data)
}

func (w *WikipediaScraperImpl) articleLinkToAbsLink(link string) string {
	/// TODO dont hardcode this.
	baseLink := w.getBaseLink("en")
	return baseLink + link
}
