package shared_test

import (
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/AppleGamer22/rake/shared"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	filePathRegularExpression = regexp.MustCompile(`\.(jpg)|(webp)|(mp4)|(webm)`)
)

func init() {
	viper.SetConfigName(".rake")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("..")
}

func TestInstagramPublicSingleImage(t *testing.T) {
	err := viper.ReadInConfig()
	assert.NoError(t, err)
	raker, err := shared.NewRaker(viper.GetString("test"), false, false)
	assert.NoError(t, err)
	URLs, username, err := raker.Instagram("CbgDyqkFBdj")
	assert.NoError(t, err)
	assert.Equal(t, "wikipedia", username)
	assert.Len(t, URLs, 1)
	URL, err := url.Parse(URLs[0])
	assert.NoError(t, err)
	assert.Equal(t, "https", URL.Scheme)
	assert.True(t, strings.Contains(URL.Host, "cdninstagram.com"))
	assert.Regexp(t, filePathRegularExpression, URL.Path)
}

func TestInstagramPublicSingleVideo(t *testing.T) {
	err := viper.ReadInConfig()
	assert.NoError(t, err)
	raker, err := shared.NewRaker(viper.GetString("test"), false, false)
	assert.NoError(t, err)
	URLs, username, err := raker.Instagram("BKyN0E2AApX")
	assert.NoError(t, err)
	assert.Equal(t, "wikipedia", username)
	assert.Len(t, URLs, 1)
	URL, err := url.Parse(URLs[0])
	assert.NoError(t, err)
	assert.Equal(t, "https", URL.Scheme)
	assert.True(t, strings.Contains(URL.Host, "cdninstagram.com"))
	assert.Regexp(t, filePathRegularExpression, URL.Path)
}

func TestInstagramPublicBundleImages(t *testing.T) {
	err := viper.ReadInConfig()
	assert.NoError(t, err)
	raker, err := shared.NewRaker(viper.GetString("test"), false, false)
	assert.NoError(t, err)
	URLs, username, err := raker.Instagram("CZNJeAil1BC")
	assert.NoError(t, err)
	assert.Equal(t, "wikipedia", username)
	assert.Len(t, URLs, 2)
	for _, urlString := range URLs {
		URL, err := url.Parse(urlString)
		assert.NoError(t, err)
		assert.Equal(t, "https", URL.Scheme)
		assert.True(t, strings.Contains(URL.Host, "cdninstagram.com"))
		assert.Regexp(t, filePathRegularExpression, URL.Path)
	}
}
