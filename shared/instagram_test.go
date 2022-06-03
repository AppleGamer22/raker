package shared_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/AppleGamer22/rake/shared"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	instagramURLRegularExpression = regexp.MustCompile(`\.(jpg)|(webp)|(mp4)|(webm)`)
)

func init() {
	viper.SetConfigName(".rake")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("..")
}

func TestInstagramPublicSingleImage(t *testing.T) {
	err := viper.ReadInConfig()
	assert.NoError(t, err)
	raker, err := shared.NewRaker(shared.FindExecutablePath(), viper.GetString("test"), false, false)
	assert.NoError(t, err)
	URLs, username, err := raker.Instagram("CbgDyqkFBdj")
	assert.NoError(t, err)
	assert.Equal(t, "wikipedia", username)
	assert.Len(t, URLs, 1)
	URL := URLs[0]
	assert.True(t, strings.HasPrefix(URL, "https://"), URL)
	assert.Regexp(t, instagramURLRegularExpression, URL)
}

func TestInstagramPublicSingleVideo(t *testing.T) {
	err := viper.ReadInConfig()
	assert.NoError(t, err)
	raker, err := shared.NewRaker(shared.FindExecutablePath(), viper.GetString("test"), false, false)
	assert.NoError(t, err)
	URLs, username, err := raker.Instagram("BKyN0E2AApX")
	assert.NoError(t, err)
	assert.Equal(t, "wikipedia", username)
	assert.Len(t, URLs, 1)
	URL := URLs[0]
	assert.True(t, strings.HasPrefix(URL, "https://"), URL)
	assert.Regexp(t, instagramURLRegularExpression, URL)
}

func TestInstagramPublicBundleImages(t *testing.T) {
	err := viper.ReadInConfig()
	assert.NoError(t, err)
	raker, err := shared.NewRaker(shared.FindExecutablePath(), viper.GetString("test"), false, false)
	assert.NoError(t, err)
	URLs, username, err := raker.Instagram("CZNJeAil1BC")
	assert.NoError(t, err)
	assert.Equal(t, "wikipedia", username)
	assert.Len(t, URLs, 2)
	for _, URL := range URLs {
		assert.True(t, strings.HasPrefix(URL, "https://"), URL)
		assert.Regexp(t, instagramURLRegularExpression, URL)
	}
}
