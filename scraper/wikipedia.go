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

var subDomainMap = map[entity.SiteCode]string{
	entity.OTZIT_UK: "en",
	entity.OTZIT_US: "en",
}

func getSubdomain(siteCode entity.SiteCode) string {
	subDomain, ok := subDomainMap[siteCode]
	if !ok {
		err := util.InvalidSiteCodeErr(siteCode)
		sentry.CaptureException(err)
		panic(err)
	}
	return subDomain
}

// WikipediaScraperImpl is an implementation a scraper
// service that will scrape wikipedia
type WikipediaScraperImpl struct{}

// NewWikipediaScraperService ...
func NewWikipediaScraperService() *WikipediaScraperImpl {
	return &WikipediaScraperImpl{}
}

func (w *WikipediaScraperImpl) getBaseLink(siteCode entity.SiteCode) string {
	subdomain := getSubdomain(siteCode)
	return fmt.Sprintf("https://%s.wikipedia.org", subdomain)
}

func (w *WikipediaScraperImpl) buildRequestURL(query string, siteCode entity.SiteCode) string {
	baseLink := w.getBaseLink(siteCode)
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

func (w *WikipediaScraperImpl) getSearchResultSet(query string, siteCode entity.SiteCode) []ScrapedResult {
	c := colly.NewCollector()
	var results []ScrapedResult

	// here we are scraping the container of all of the result boxes.
	c.OnHTML("div .mw-search-results", func(e *colly.HTMLElement) {
		e.DOM.Find(".mw-search-result").Each(func(_ int, body *goquery.Selection) {
			title := body.Find(".mw-search-result-heading a").Text()
			meta := body.Find(".searchresult").Text()

			pageLink, ok := body.Find(".mw-search-result-heading a").Attr("href")
			if !ok {
				sentry.CaptureException(fmt.Errorf("failed to get URL for result '%s'", title))
				return
			}

			link := w.articleLinkToAbsLink(pageLink, siteCode)

			convertedLink, ok := w.convertLink(link)
			if !ok {
				sentry.CaptureException(fmt.Errorf("failed to convert URL for result '%s'", title))
				return
			}

			results = append(results, ScrapedResult{
				Title: title,
				Href: convertedLink,
				Snippet: meta,
			})
		})
	})

	url := w.buildRequestURL(query, siteCode)
	if err := c.Visit(url); err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	return results
}

// Scrape will scrape google for the given query and
// parse the results.
func (w *WikipediaScraperImpl) Scrape(query string, siteCode entity.SiteCode) []ScrapedResult {
	convertedResults := w.getSearchResultSet(query, siteCode)
	return convertedResults
}

func (w *WikipediaScraperImpl) articleLinkToAbsLink(link string, siteCode entity.SiteCode) string {
	baseLink := w.getBaseLink(siteCode)
	return baseLink + link
}
