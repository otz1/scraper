package scraper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItScrapesGoogleResults(t *testing.T) {
	gscraper := NewGoogleScraperService()
	results := gscraper.Scrape("how to make pancakes")
	assert.Equal(t, len(results), len(results))
}