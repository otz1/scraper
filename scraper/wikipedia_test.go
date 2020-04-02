package scraper

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestWikipediaQueryReturnsSomeResults(t *testing.T) {
	wss := NewWikipediaScraperService()
	results := wss.Scrape("how to make pancakes")
	assert.NotEqual(t, 0, len(results))
}

// valid meaning we have a protocol at least.
func TestWikipediaQueryContainsValidLinks(t *testing.T) {
	wss := NewWikipediaScraperService()
	results := wss.Scrape("how to make pancakes")
	for _, result := range results {
		_, err := url.Parse(result.Href)
		assert.NoError(t, err)
	}
	assert.NotEqual(t, 0, len(results))
}