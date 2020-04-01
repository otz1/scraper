package entity

import "github.com/otz1/scraper/scraper"

type Result struct {
	Title string
	Href string
}

func ToResult(sr scraper.ScrapedResult) Result {
	return Result{
		Title: sr.Title,
		Href: sr.Href,
	}
}

type ScrapeResponse struct {
	OriginalQuery string
	Results []Result
}
