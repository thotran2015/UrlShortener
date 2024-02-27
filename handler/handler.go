package handler

import (
	"UrlShortener/shortener"
	"UrlShortener/store"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync/atomic"
)

const HOST = "http://localhost:9808/"

// User ID incrementor
var userIdCounter int64

// Expose the service functionalities at two API endpoints: shorten URL and redirect to original URL

type UrlCreateRequest struct {
	OriUrl string `json:"long_url" binding:"required"`
}

func CreateShortUrl(c *gin.Context) {
	var createRequest UrlCreateRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation error": err.Error()})
		return
	}
	fmt.Printf("Long URL: %s\n", createRequest.OriUrl)
	shortUrl := shortener.GenerateShortLink(createRequest.OriUrl, strconv.FormatInt(userIdCounter, 10))
	fmt.Printf("Short URL: %s %s\n", shortUrl, HOST+shortUrl)
	store.SaveUrlMapping(shortUrl, createRequest.OriUrl)
	atomic.AddInt64(&userIdCounter, 1)
	c.JSON(200, gin.H{
		"message":   "short url created successfully",
		"short_url": HOST + shortUrl,
	})
}

func RedirectShortUrl(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	originalUrl := store.GetUrlMapping(shortUrl)
	// use redirect func from gin context
	c.Redirect(302, originalUrl)
}
