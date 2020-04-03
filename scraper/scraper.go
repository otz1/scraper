package scraper

import (
	"github.com/otz1/scraper/entity"
	"log"
	"net/url"
	"strings"
)

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

func (sr ScrapedResult) ToResult() entity.Result {
	return entity.Result{
		Title: sr.Title,
		Href: sr.Href,
		Snippet: sr.Snippet,
	}
}

func attrsToMap(attrs string) map[string]bool {
	result := map[string]bool{}
	for _, a := range strings.Split(attrs, " ") {
		result[a] = true
	}
	return result
}

// FIXME dead code?
func elHasClass(attribs string, class string) bool {
	classes := attrsToMap(attribs)
	log.Println("is", class, "inside of", classes)
	_, ok := classes[class]
	return ok
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
