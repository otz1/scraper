package resource

import (
	"fmt"

	"github.com/otz1/scraper/entity"
	"github.com/otz1/scraper/scraper"
)

// GoogleSearchResource implements the scraping for
// google as a source of information
type GoogleSearchResource struct{}

// NewGoogleSearchResource ...
func NewGoogleSearchResource() *GoogleSearchResource {
	return &GoogleSearchResource{}
}

// Query ...
func (gss *GoogleSearchResource) Query(query string) entity.ScrapeResponse {
	gscraper := scraper.NewGoogleScraperService()
	results := gscraper.Scrape(query)
	fmt.Println("results are ", results)

	// todo store results in response and convert where necessary
	return entity.ScrapeResponse{
		OriginalQuery: query,
	}
}
