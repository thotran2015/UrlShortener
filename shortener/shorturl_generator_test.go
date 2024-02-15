package shortener

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateShortLink(t *testing.T) {
	originalUrl1 := "https://thotran2015.github.io/portfolio018/"
	shortUrl1 := GenerateShortLink(originalUrl1, "0")
	assert.Equal(t, "6zKgqaXy", shortUrl1)
	assert.Equal(t, UrlLength, len(shortUrl1))
	shortUrl2 := GenerateShortLink(originalUrl1, "1")
	assert.Equal(t, "hV2MJc13", shortUrl2)
	assert.Equal(t, UrlLength, len(shortUrl2))
}
