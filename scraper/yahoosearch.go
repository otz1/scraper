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

var yahooSiteMap = map[entity.SiteCode]string{
	entity.OTZIT_US: "uk.search",
	entity.OTZIT_UK: "uk.search",
}

func (y *YahooSearchScraperImpl) getSubdomain(siteCode entity.SiteCode) string {
	sub, ok := yahooSiteMap[siteCode]
	if !ok {
		err := util.InvalidSiteCodeErr(siteCode)
		sentry.CaptureException(err)
		panic(err)
	}
	return sub
}

// YahooSearchScraperImpl is an implementation a scraper
// service that will scrape yahoo
type YahooSearchScraperImpl struct {
	BasicScraper
}

// NewYahooScraperService ...
func NewYahooScraperService() *YahooSearchScraperImpl {
	return &YahooSearchScraperImpl{
		NewBasicScraper(),
	}
}

func (y *YahooSearchScraperImpl) buildRequestURL(query string, siteCode entity.SiteCode) string {
	sub := y.getSubdomain(siteCode)
	return fmt.Sprintf(`https://%s.yahoo.com/search?p=%s`, sub, parseQuery(query))
}

func (y *YahooSearchScraperImpl) convertLink(link string) (string, bool) {
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

func (y *YahooSearchScraperImpl) getSearchResultSet(query string, siteCode entity.SiteCode) []ScrapedResult {
	var results []ScrapedResult

	// here we are scraping the container of all of the result boxes.
	y.collector.OnHTML("div[id=left]", func(e *colly.HTMLElement) {
		e.DOM.Find("div.dd.Sr").Each(func(_ int, body *goquery.Selection) {
			title := body.Find(".compTitle h3").Text()
			snippet := body.Find(".compText p").Text()

			link, ok := body.Find(".compTitle h3 a").Attr("href")
			if !ok {
				sentry.CaptureException(fmt.Errorf("failed to get URL for result '%s'", title))
				return
			}

			convertedLink, ok := y.convertLink(link)
			if !ok {
				sentry.CaptureException(fmt.Errorf("failed to convert URL for result '%s'", title))
				return
			}

			results = append(results, ScrapedResult{
				Title:   title,
				Href:    convertedLink,
				Snippet: snippet,
			})
		})
	})

	url := y.buildRequestURL(query, siteCode)
	if err := y.collector.Visit(url); err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	return results
}

// Scrape will scrape google for the given query and
// parse the results.
func (y *YahooSearchScraperImpl) Scrape(query string, siteCode entity.SiteCode) []ScrapedResult {
	convertedResults := y.getSearchResultSet(query, siteCode)
	return convertedResults
}
