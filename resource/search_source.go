package resource

import (
	"github.com/otz1/scraper/entity"
)

// SearchResource will handle the endpoint for
// a given search service that can be scraped
type SearchResource interface {
	Query(query string, siteCode entity.SiteCode) entity.ScrapeResponse
}

// ValidSearchResources is a list of valid resources that can be used
// for processing scrapes.
func ValidSearchResources() map[entity.ScrapeSource]SearchResource {
	return map[entity.ScrapeSource]SearchResource{
		entity.DDG:       NewDDGSearchResource(),
		entity.WIKIPEDIA: NewWikipediaSearchResource(),
	}
}
