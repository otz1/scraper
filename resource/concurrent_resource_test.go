package resource

import (
	"github.com/otz1/scraper/entity"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestItScrapesAllAvailableSources(t *testing.T) {
	resp := ScrapeAvailableSources("how to make pancakes", entity.OTZIT_UK)
	assert.NotEmpty(t, resp.Results)

	availableSources := ValidSearchResources()

	scrapedSources := getSources(resp.Results)
	for source := range availableSources {
		if _, ok := scrapedSources[source]; !ok {
			assert.Fail(t, "didn't scrape all sources!")
		}
	}

	log.Println("we scraped", scrapedSources)
}

func getSources(results []entity.Result) map[entity.ScrapeSource]bool {
	r := map[entity.ScrapeSource]bool{}
	for _, res := range results {
		if _, ok := r[res.Source]; !ok {
			r[res.Source] = true
		}
	}
	return r
}
