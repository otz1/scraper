package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/otz1/scraper/entity"
	"github.com/otz1/scraper/resource"
)

// TODO move this somewhere else.
func ScrapeHandler(c *gin.Context) {
	var req entity.ScrapeRequest
	if err := c.BindJSON(&req); err != nil {
		panic(err)
	}

	selectedSource := GOOGLE
	if source := req.Source; source == nil {
		selectedSource = *source
	}

	sources := resource.ValidSearchResources()
	source, ok := sources[source]
	if !ok {
		panic("bad source!")
	}

	// do the scrape on the given source.
	source.Scrape(req.Query)
}

func main() {
	router := gin.Default()
	{
		conf := cors.Default()
		router.Use(conf)
	}
	router.GET("/scrape", ScrapeHandler)
	router.Run(":8001")
}
