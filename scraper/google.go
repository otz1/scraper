package scraper

// https://edmundmartin.com/scraping-google-with-golang/

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/otz1/scraper/entity"
	"github.com/otz1/scraper/util"
	"log"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

// TODO migrate over to using sitecode.

var siteCodeMap = map[entity.SiteCode]string{
	entity.OTZIT_US: "com",
	entity.OTZIT_UK: "co.uk",
	entity.OTZIT_FR: "fr",
}

func (g GoogleScraperImpl) getTLD(siteCode entity.SiteCode) string {
	tld, ok := siteCodeMap[siteCode]
	if !ok {
		err := util.InvalidSiteCodeErr(siteCode)
		sentry.CaptureException(err)
		panic(err)
	}
	return tld
}

// GoogleScraperImpl is an impl. of a scraper
// service that will scrape google
type GoogleScraperImpl struct{}

// NewGoogleScraperService ...
func NewGoogleScraperService() *GoogleScraperImpl {
	return &GoogleScraperImpl{}
}

func (g *GoogleScraperImpl) buildRequestURL(query string, siteCode entity.SiteCode) string {
	tld := g.getTLD(siteCode)
	return fmt.Sprintf(`https://www.google.%s/search?q=%s&hl=%s`, tld, parseQuery(query), tld)
}

func (g *GoogleScraperImpl) getSearchResultSet(query string, siteCode entity.SiteCode) []*colly.HTMLElement {
	c := colly.NewCollector()
	var searchElements []*colly.HTMLElement
	c.OnHTML("div", func(e *colly.HTMLElement) {

		e.ForEach("#search", func(count int, e *colly.HTMLElement) {
			log.Println(e)
		})

		// searchElements = append(searchElements, e)
	})
	url := g.buildRequestURL(query, siteCode)
	if err := c.Visit(url); err != nil {
		sentry.CaptureException(err)
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

			// titles := e.ChildTexts("h3")

			// ensure it's a proper link!
			if _, err := url.Parse(link); err != nil {
				sentry.CaptureException(err)
				return
			}

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
func (g *GoogleScraperImpl) Scrape(query string, siteCode entity.SiteCode) []ScrapedResult {
	resultSet := g.getSearchResultSet(query, siteCode)
	convertedResults := g.convertResults(resultSet)
	return convertedResults
}
