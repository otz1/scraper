package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/otz1/scraper/scrapecache"
	"github.com/otz1/scraper/util"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/otz1/scraper/entity"
)

// we cache the entire response for a minute.
// in the future we could just cache results. for now this works.
var cachedScraper = scrapecache.New()

func ScrapeHandler(c *gin.Context) {
	var req entity.ScrapeRequest
	if err := c.BindJSON(&req); err != nil {
		sentry.CaptureException(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	siteCode := util.GetSiteCode(c.GetHeader("SITE-CODE"))
	resp := cachedScraper.Query(siteCode, req.Query)
	c.JSON(http.StatusOK, resp)
}

func main() {
	env := util.GetEnv("ENVIRONMENT", "local")
	log.Println("RUNNING IN", strings.ToUpper(env))

	err := sentry.Init(sentry.ClientOptions{
		Dsn:         "https://076afa24ea2b4cdd904ff677b5f92f62@sentry.io/5187016",
		Environment: env,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	router := gin.Default()
	{
		conf := cors.Default()
		router.Use(conf)
	}
	router.POST("/scrape", ScrapeHandler)
	router.Run(":8001")
}
