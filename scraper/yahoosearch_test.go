package scraper

import (
	"github.com/otz1/scraper/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItScrapesYahooSearches(t *testing.T) {
	yscraper := NewYahooScraperService()
	resp := yscraper.Scrape("how to make a cup of tea", entity.OTZIT_UK)
	assert.NotEmpty(t, resp)
}
