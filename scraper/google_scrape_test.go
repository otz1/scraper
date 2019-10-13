package scraper

import (
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/assert"
)

func TestItBuildsTheURL(t *testing.T) {

}

func TestItScrapesGoogle(t *testing.T) {
	query := randomdata.SillyName()
	s := NewGoogleScraperService()
	results := s.Scrape(query)
	assert.NotEmpty(t, results)
}
