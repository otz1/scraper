package scraper

import "net/url"

type ScrapedResult struct {
	Title   string
	Href    string
	Keyword []string
}

type ScraperService interface {
	Scrape(query string) []ScrapedResult
}

func parseQuery(query string) string {
	return url.PathEscape(query)
}
