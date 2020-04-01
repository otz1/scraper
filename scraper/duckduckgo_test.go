package scraper

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestQueryReturnsSomeResults(t *testing.T) {
	ddg := NewDDGScraperService()
	results := ddg.Scrape("how to make pancakes")
	assert.NotEqual(t, 0, len(results))
}

// valid meaning we have a protocol at least.
func TestQueryContainsValidLinks(t *testing.T) {
	ddg := NewDDGScraperService()
	results := ddg.Scrape("how to make pancakes")
	for _, result := range results {
		_, err := url.Parse(result.Href)
		assert.NoError(t, err)
	}
	assert.NotEqual(t, 0, len(results))
}