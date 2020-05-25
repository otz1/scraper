package resource

import (
	"github.com/otz1/scraper/entity"
	"sync"
)

var resources = ValidSearchResources()

// ScrapeAvailableSources will concurrently scrape from all resources available.
func ScrapeAvailableSources(query string, siteCode entity.SiteCode) entity.ScrapeResponse {
	var wg sync.WaitGroup
	wg.Add(len(resources))

	// we communicate between routines with a channel that can
	// buffer one response at a time.
	queue := make(chan entity.ScrapeResponse, 1)

	// this is where we do the queries and
	// queue the responses in multiple routines
	for _, res := range resources {
		go func(res SearchResource) {
			queue <- res.Query(query, siteCode)
		}(res)
	}

	// we have one routine that exists throughout the entire scrape
	// which will take all the queued items and put them in a single
	// store.
	var scrapedResults []entity.Result
	go func() {
		for result := range queue {
			scrapedResults = append(scrapedResults, result.Results...)
			wg.Done()
		}
	}()
	wg.Wait()

	return entity.ScrapeResponse{
		OriginalQuery: query,
		Results:       scrapedResults,
	}
}