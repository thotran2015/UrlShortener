package main

import (
	"UrlShortener/handler"
	"UrlShortener/store"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default() // Create a router, gin engine
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "This is a Go URL Shortener!"})
	})
	// Create a short URL with POST
	router.POST("/create-short-url", func(c *gin.Context) {
		handler.CreateShortUrl(c)
	})
	// Redirect short URL to original one
	// :shortUrl indicate path parameter called shortUrl
	router.GET("/:shortUrl", func(c *gin.Context) {
		handler.RedirectShortUrl(c)
	})

	// Initialize store before starting URL shortening service
	store.InitializeStore()

	err := router.Run(":9808")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server. Error: %v", err))
	}
}
