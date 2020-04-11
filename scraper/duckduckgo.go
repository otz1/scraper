package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/getsentry/sentry-go"
	"github.com/gocolly/colly"
	"github.com/otz1/scraper/entity"
	"github.com/otz1/scraper/util"
	"log"
	"net/url"
	"strings"
)

var siteEnumMap = map[entity.SiteCode]string {
	entity.OTZIT_US: "com",
	entity.OTZIT_UK: "co.uk",
}

func getTLD(siteCode entity.SiteCode) string {
	tld, ok := siteEnumMap[siteCode]
	if !ok {
		err := util.InvalidSiteCodeErr(siteCode)
		sentry.CaptureException(err)
		panic(err)
	}
	return tld
}

// DuckDuckGoScraperImpl is an implementation a scraper
// service that will scrape duckduckgo
type DuckDuckGoScraperImpl struct{}

// NewDDGScraperService ...
func NewDDGScraperService() *DuckDuckGoScraperImpl {
	return &DuckDuckGoScraperImpl{}
}

func (d *DuckDuckGoScraperImpl) buildRequestURL(query string, siteCode entity.SiteCode) string {
	tld := getTLD(siteCode)
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

func (d *DuckDuckGoScraperImpl) getSearchResultSet(query string, siteCode entity.SiteCode) []ScrapedResult {
	c := colly.NewCollector()
	var results []ScrapedResult

	// here we are scraping the container of all of the result boxes.
	c.OnHTML("div[id=links]", func(e *colly.HTMLElement) {
		e.DOM.Find(".result").Each(func(_ int, s *goquery.Selection) {
			body := s.Find(".result__body")
			title := body.Find(".result__title a").Text()
			snippet := body.Find(".result__snippet").Text()

			link, ok := body.Find(".result__title .result__a").Attr("href")
			if !ok {
				sentry.CaptureException(fmt.Errorf("failed to get URL for result '%s'", title))
				return
			}

			convertedLink, ok := d.convertLink(link)
			if !ok {
				sentry.CaptureException(fmt.Errorf("failed to convert URL for result '%s'", title))
				return
			}

			results = append(results, ScrapedResult{
				Title: title,
				Href: convertedLink,
				Snippet: snippet,
			})
		})
	})

	url := d.buildRequestURL(query, siteCode)
	if err := c.Visit(url); err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	return results
}

// Scrape will scrape google for the given query and
// parse the results.
func (d *DuckDuckGoScraperImpl) Scrape(query string, siteCode entity.SiteCode) []ScrapedResult {
	convertedResults := d.getSearchResultSet(query, siteCode)
	return convertedResults
}