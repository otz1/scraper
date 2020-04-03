package scraper

import (
	"github.com/otz1/scraper/entity"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestDDGQueryReturnsSomeResults(t *testing.T) {
	ddg := NewDDGScraperService()
	results := ddg.Scrape("how to make pancakes", entity.OTZIT_UK)
	assert.NotEqual(t, 0, len(results))
}

// valid meaning we have a protocol at least.
func TestDDGQueryContainsValidLinks(t *testing.T) {
	ddg := NewDDGScraperService()
	results := ddg.Scrape("how to make pancakes", entity.OTZIT_UK)
	for _, result := range results {
		_, err := url.Parse(result.Href)
		assert.NoError(t, err)
	}
	assert.NotEqual(t, 0, len(results))
}