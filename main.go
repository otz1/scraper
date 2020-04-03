package main

import (
	"github.com/getsentry/sentry-go"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/otz1/scraper/entity"
	"github.com/otz1/scraper/resource"
)

// TODO move this somewhere else.
func ScrapeHandler(c *gin.Context) {
	var req entity.ScrapeRequest
	if err := c.BindJSON(&req); err != nil {
		log.Println("bad request", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	selectedSource := entity.DDG
	if source := req.Source; source != nil {
		selectedSource = *source
	}

	sources := resource.ValidSearchResources()
	resource, ok := sources[selectedSource]
	if !ok {
		panic("bad source!")
	}

	// do the scrape on the given source.
	resp := resource.Query(req.Query)
	c.JSON(http.StatusOK, resp)
}

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://076afa24ea2b4cdd904ff677b5f92f62@sentry.io/5187016",
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
