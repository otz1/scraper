package scraper

import "net/url"

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

// Service is a service that scrapes
// for the given query from a particular source
type Service interface {
	Scrape(query string) []ScrapedResult
}

func parseQuery(query string) string {
	return url.PathEscape(query)
}
