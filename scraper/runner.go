package scraper

import "log"

// ScrapeFor will scrape for a query on all
// the services available. it will handle
// query parsing too
func ScrapeFor(query string) {
	scrapers := []ScraperService{
		NewGoogleScraperService(),
	}

	log.Println("We have", len(scrapers), "plumbed in")
	log.Println("->", scrapers)

	// TODO make this parallel
	for _, ss := range scrapers {
		ss.Scrape(parseQuery(query))
	}
}
