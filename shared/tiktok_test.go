package shared_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/AppleGamer22/rake/shared"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	viper.SetConfigName(".rake")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("..")
}

func TestTikTokPublicSingleVideo(t *testing.T) {
	raker, err := shared.NewRaker("", false, true)
	assert.NoError(t, err)
	urlString, username, err := raker.TikTok("f1", "7048983181063687430")
	assert.NoError(t, err)
	assert.Equal(t, "f1", username)
	URL, err := url.Parse(urlString)
	assert.NoError(t, err)
	assert.Equal(t, "https", URL.Scheme)
	assert.True(t, strings.Contains(URL.Host, "-webapp.tiktok.com"), URL)
}
