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
func (ddg *DDGSearchResource) Query(query string, siteCode entity.SiteCode) entity.ScrapeResponse {
	ddgscraper := scraper.NewDDGScraperService()
	scrapedResults := ddgscraper.Scrape(query, siteCode)

	convertedResults := func() []entity.Result {
		var results []entity.Result
		for _, sr := range scrapedResults {
			result := sr.ToResult()
			results = append(results, result)
		}
		return results
	}()

	return entity.ScrapeResponse{
		OriginalQuery: query,
		Results: convertedResults,
	}
}