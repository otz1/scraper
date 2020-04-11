package scraper

import (
	"github.com/otz1/scraper/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItScrapesGoogleResults(t *testing.T) {
	gscraper := NewGoogleScraperService()
	results := gscraper.Scrape("how to make pancakes", entity.OTZIT_UK)
	assert.Equal(t, len(results), len(results))
}
