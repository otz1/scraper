package entity

import "github.com/otz1/scraper/scraper"

type Result struct {
	Title string
	Href string
	Snippet string
}

func ToResult(sr scraper.ScrapedResult) Result {
	return Result{
		Title: sr.Title,
		Href: sr.Href,
		Snippet: sr.Snippet,
	}
}

type ScrapeResponse struct {
	OriginalQuery string
	Results []Result
}
