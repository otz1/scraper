package resource

import (
	"github.com/otz1/scraper/entity"
	"github.com/otz1/scraper/scraper"
)

// WikipediaSearchResource implements the scraping for
// wikipedia as a source of information
type WikipediaSearchResource struct {
	scraperService *scraper.WikipediaScraperImpl
}

// NewGoogleSearchResource ...
func NewWikipediaSearchResource() *WikipediaSearchResource {
	return &WikipediaSearchResource{
		scraperService: scraper.NewWikipediaScraperService(),
	}
}

// Query ...
func (wsr *WikipediaSearchResource) Query(query string, siteCode entity.SiteCode) entity.ScrapeResponse {
	scrapedResults := wsr.scraperService.Scrape(query, siteCode)

	convertedResults := make([]entity.Result, len(scrapedResults))
	for i, sr := range scrapedResults {
		result := sr.ToResult(entity.WIKIPEDIA)
		convertedResults[i] = result
	}

	return entity.ScrapeResponse{
		OriginalQuery: query,
		Results:       convertedResults,
	}
}
