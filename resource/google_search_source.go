package resource

import (
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
	scrapedResults := gscraper.Scrape(query)

	convertedResults := func() []entity.Result {
		var results []entity.Result
		for _, sr := range scrapedResults {
			result := entity.ToResult(sr)
			results = append(results, result)
		}
		return results
	}()

	return entity.ScrapeResponse{
		OriginalQuery: query,
		Results: convertedResults,
	}
}
