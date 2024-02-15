package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testStoreService = &StorageService{}

func init() {
	testStoreService = InitializeStore()
}

func TestStoreInit(t *testing.T) {
	assert.True(t, testStoreService.redisClient != nil)
}

func TestSetAndGet(t *testing.T) {
	originalUrl := "https://thotran2015.github.io/portfolio018/"
	shortUrl := "https://th018/"

	SaveUrlMapping(shortUrl, originalUrl)
	longUrl := GetUrlMapping(shortUrl)
	assert.Equal(t, originalUrl, longUrl)
}
