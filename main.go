package main

import (
	"log"
	"net/http"

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
	router := gin.Default()
	{
		conf := cors.Default()
		router.Use(conf)
	}
	router.POST("/scrape", ScrapeHandler)
	router.Run(":8001")
}
