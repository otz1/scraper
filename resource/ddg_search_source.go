package resource

import (
	"github.com/otz1/scraper/entity"
	"github.com/otz1/scraper/scraper"
)

// DDGSearchResource implements the scraping for
// duckduckgo as a source of information
type DDGSearchResource struct{}

// NewGoogleSearchResource ...
func NewDDGSearchResource() *DDGSearchResource {
	return &DDGSearchResource{}
}

// Query ...
func (ddg *DDGSearchResource) Query(query string) entity.ScrapeResponse {
	ddgscraper := scraper.NewDDGScraperService()
	scrapedResults := ddgscraper.Scrape(query)

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