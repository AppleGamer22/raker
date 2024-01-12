package shared_test

import (
	"net/url"
	"regexp"
	"testing"

	"github.com/AppleGamer22/raker/shared"
	"github.com/stretchr/testify/assert"
)

var filePathRegularExpression = regexp.MustCompile(`\.(jpg)|(webp)|(mp4)|(webm)`)
var instagramDomainRegularExpression = regexp.MustCompile(`(cdninstagram\.com)|(fbcdn\.net)`)

func testInstagramURLs(t *testing.T, URLs []string) {
	for _, urlString := range URLs {
		URL, err := url.Parse(urlString)
		assert.NoError(t, err)
		assert.Equal(t, "https", URL.Scheme)
		assert.Regexp(t, instagramDomainRegularExpression, URL.Host, urlString)
		assert.Regexp(t, filePathRegularExpression, URL.Path)
	}
}

func TestInstagramSingleImage(t *testing.T) {
	instagram := shared.NewInstagram(configuration.Instagram.FBSR, configuration.Instagram.Session, configuration.Instagram.User)
	URLs, username, err := instagram.Post("CbgDyqkFBdj")
	assert.NoError(t, err)
	assert.Equal(t, "wikipedia", username)
	assert.Len(t, URLs, 1)
	testInstagramURLs(t, URLs)
}

func TestInstagramSingleVideo(t *testing.T) {
	instagram := shared.NewInstagram(configuration.Instagram.FBSR, configuration.Instagram.Session, configuration.Instagram.User)
	URLs, username, err := instagram.Post("BKyN0E2AApX")
	assert.NoError(t, err)
	assert.Equal(t, "wikipedia", username)
	assert.Len(t, URLs, 1)
	testInstagramURLs(t, URLs)
}

func TestInstagramBundledImages(t *testing.T) {
	instagram := shared.NewInstagram(configuration.Instagram.FBSR, configuration.Instagram.Session, configuration.Instagram.User)
	URLs, username, err := instagram.Post("CZNJeAil1BC")
	assert.NoError(t, err)
	assert.Equal(t, "wikipedia", username)
	assert.Len(t, URLs, 2)
	testInstagramURLs(t, URLs)
}
