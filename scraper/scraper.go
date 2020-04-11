package scraper

import (
	"github.com/gocolly/colly"
	"github.com/otz1/scraper/entity"
	"net/url"
)

type BasicScraper struct {
	collector *colly.Collector
}

func NewBasicScraper(options ...colly.CollectorOption) BasicScraper {
	options = append(options, colly.AllowURLRevisit())
	return BasicScraper{
		collector: colly.NewCollector(options...),
	}
}

// ScrapedResult is a result that has been scraped from
// one of the external sources.
type ScrapedResult struct {
	Title           string `json:"title"`
	Snippet         string `json:"snippet"`
	Ranking         int    `json:"ranking"`
	ImageSource     string `json:"imageSource"`
	ThumbnailSource string `json:"thumbnailSource"`
	Href            string `json:"href"`
}

// This should go in conv.
func (sr ScrapedResult) ToResult(source entity.ScrapeSource) entity.Result {
	return entity.Result{
		Title:   sr.Title,
		Href:    sr.Href,
		Snippet: sr.Snippet,
		Source:  source,
	}
}

// Service is a service that scrapes
// for the given query from a particular source
type Service interface {
	Scrape(query string, siteCode entity.SiteCode) []ScrapedResult
	buildRequestURL(query string, siteCode entity.SiteCode) string
}

func parseQuery(query string) string {
	return url.PathEscape(query)
}
