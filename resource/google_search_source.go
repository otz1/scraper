package resource

import "github.com/otz1/scraper/scraper"

// GoogleSearchResource implements the scraping for
// google as a source of information
type GoogleSearchResource struct{}

// NewGoogleSearchResource ...
func NewGoogleSearchResource() GoogleSearchResource {
	return GoogleSearchResource{}
}

// Query ...
func (gss *GoogleSearchResource) Query(query string) {
	gscraper := scraper.NewGoogleScraperService()
	results := gscraper.Scrape(query)
	panic("do some stuff here")
}
