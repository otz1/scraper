package resource

import (
	"github.com/otz1/scraper/entity"
	"github.com/otz1/scraper/scraper"
)

// YahooSearchResource implements the scraping for
// wikipedia as a source of information
type YahooSearchResource struct {
	scraperService *scraper.YahooSearchScraperImpl
}

// NewYahooSearchResource ...
func NewYahooSearchResource() *YahooSearchResource {
	return &YahooSearchResource{
		scraperService: scraper.NewYahooScraperService(),
	}
}

// Query ...
func (ysr *YahooSearchResource) Query(query string, siteCode entity.SiteCode) entity.ScrapeResponse {
	scrapedResults := ysr.scraperService.Scrape(query, siteCode)

	convertedResults := make([]entity.Result, len(scrapedResults))
	for i, sr := range scrapedResults {
		result := sr.ToResult(entity.YAHOO)
		convertedResults[i] = result
	}

	return entity.ScrapeResponse{
		OriginalQuery: query,
		Results:       convertedResults,
	}
}
