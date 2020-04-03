package resource

import (
	"github.com/otz1/scraper/entity"
	"github.com/otz1/scraper/scraper"
)

// WikipediaSearchResource implements the scraping for
// wikipedia as a source of information
type WikipediaSearchResource struct{}

// NewGoogleSearchResource ...
func NewWikipediaSearchResource() *WikipediaSearchResource {
	return &WikipediaSearchResource{}
}

// Query ...
func (wsr *WikipediaSearchResource) Query(query string, siteCode entity.SiteCode) entity.ScrapeResponse {
	wscraper := scraper.NewWikipediaScraperService()
	scrapedResults := wscraper.Scrape(query, siteCode)

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
