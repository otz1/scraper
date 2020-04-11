package scrapecache

import (
	"github.com/otz1/scraper/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItCachesCorrectly(t *testing.T) {
	scraper := New()

	// first cache
	resp := scraper.Query(entity.OTZIT_UK, "how to make pancakes")
	assert.NotEmpty(t, resp.Results)
	assert.Equal(t, uint64(1), scraper.caches)

	// second cache with a hit.
	resp2 := scraper.Query(entity.OTZIT_UK, "how to make pancakes")
	assert.Equal(t, resp.Results, resp2.Results)
	assert.Equal(t, uint64(1), scraper.hits)
}

// in this case we are testing that different sources and sitecodes are separated
// properly, so we dont cache UK results with FR results for example
func TestItHashesCachesCorrectly(t *testing.T) {
	scraper := New()

	// Given we have the same queries,
	// same engines, different countries.

	// When we scrape once
	resp := scraper.Query(entity.OTZIT_UK, "how to make pancakes")
	assert.NotEmpty(t, resp.Results)
	assert.Equal(t, uint64(1), scraper.caches)

	// Then we scrape the same query again from a different site
	resp2 := scraper.Query(entity.OTZIT_US, "how to make pancakes")
	assert.NotEmpty(t, resp2.Results)

	// Then we have zero hits and 2 separate caches.
	assert.Equal(t, uint64(0), scraper.hits)
	assert.Equal(t, uint64(2), scraper.caches)
}
