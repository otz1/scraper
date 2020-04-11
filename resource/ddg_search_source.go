package resource

import (
	"github.com/otz1/scraper/entity"
	"github.com/otz1/scraper/scraper"
)

// DDGSearchResource implements the scraping for
// duckduckgo as a source of information
type DDGSearchResource struct {
	scraperService *scraper.DuckDuckGoScraperImpl
}

// NewGoogleSearchResource ...
func NewDDGSearchResource() *DDGSearchResource {
	return &DDGSearchResource{
		scraperService: scraper.NewDDGScraperService(),
	}
}

// Query ...
func (ddg *DDGSearchResource) Query(query string, siteCode entity.SiteCode) entity.ScrapeResponse {
	scrapedResults := ddg.scraperService.Scrape(query, siteCode)

	// TODO move to conv
	convertedResults := make([]entity.Result, len(scrapedResults))
	for i, sr := range scrapedResults {
		result := sr.ToResult(entity.DDG)
		convertedResults[i] = result
	}

	return entity.ScrapeResponse{
		OriginalQuery: query,
		Results:       convertedResults,
	}
}
