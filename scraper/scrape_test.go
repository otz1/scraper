package scraper

import "testing"

func TestItScrapesWithAllServices(t *testing.T) {
	ScrapeFor("hello, world!")
}
