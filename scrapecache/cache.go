package scrapecache

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	jsoniter "github.com/json-iterator/go"
	"github.com/otz1/scraper/entity"
	"github.com/otz1/scraper/resource"
	"github.com/patrickmn/go-cache"
	"time"
)

// ScrapeCache is a shim over the scraping resources
// which will cache the responses over a short period of time.
type ScrapeCache struct {
	store *cache.Cache
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

	rawCachedResp, found := c.store.Get(key)
	if !found {
		c.misses++
	}

	if found {
		// we've found it, let's unmarshal
		var cachedResp entity.ScrapeResponse
		err := jsoniter.Unmarshal([]byte(rawCachedResp.(string)), &cachedResp)
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
		c.store.Set(key, string(respJSON), cache.DefaultExpiration)
	}

	return resp
}


func New() *ScrapeCache {
	return &ScrapeCache{
		store: cache.New(60*time.Second, 5*time.Minute),
		caches: 0,
		misses: 0,
		failures: 0,
		hits: 0,
	}
}
