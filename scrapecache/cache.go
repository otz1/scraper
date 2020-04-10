package scrapecache

import (
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/getsentry/sentry-go"
	jsoniter "github.com/json-iterator/go"
	"github.com/otz1/scraper/entity"
	"github.com/otz1/scraper/resource"
	"time"
)

// ScrapeCache is a shim over the scraping resources
// which will cache the responses over a short period of time.
type ScrapeCache struct {
	store *bigcache.BigCache
	caches uint64
	hits uint64
	misses uint64
	failures uint64
}

func hash(siteCode entity.SiteCode, selectedSource entity.ScrapeSource, query string) string {
	return fmt.Sprintf("%s:%s:%s", string(siteCode), string(selectedSource), query)
}

func (c *ScrapeCache) Query(siteCode entity.SiteCode, selectedSource entity.ScrapeSource, query string) entity.ScrapeResponse {
	key := hash(siteCode, selectedSource, query)

	rawCachedResp, err := c.store.Get(key)
	if err != nil {
		c.misses++
		sentry.CaptureException(err)
	}

	if err == nil {
		// we've found it, let's unmarshal
		var cachedResp entity.ScrapeResponse
		err := jsoniter.Unmarshal(rawCachedResp, &cachedResp)
		if err == nil {
			c.hits++
			return cachedResp
		}

		// cant unmarshal so we treat it as a cache miss.
		sentry.CaptureException(err)
		c.failures++
	}

	// at this point we've either
	// 1. not found the cached item
	// 2. failed to unmarshal due to potential corruption
	// so in this case we honour the scrape request.

	sources := resource.ValidSearchResources()
	scraperResource, ok := sources[selectedSource]
	if !ok {
		// TODO log as sentry error.
		panic("bad source!")
	}

	// do the scrape on the given source.
	resp := scraperResource.Query(query, siteCode)

	// store in the cache
	{
		respJSON, err := jsoniter.Marshal(&resp)
		if err != nil {
			sentry.CaptureException(err)
		}
		c.caches++
		c.store.Set(key, respJSON)
	}

	return resp
}


func New() *ScrapeCache {
	store, _ := bigcache.NewBigCache(bigcache.DefaultConfig(5 * time.Minute))
	return &ScrapeCache{
		store: store,
		caches: 0,
		misses: 0,
		failures: 0,
		hits: 0,
	}
}
