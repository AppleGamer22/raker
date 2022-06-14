package shared_test

import (
	"net/url"
	"testing"

	"github.com/AppleGamer22/rake/shared"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	viper.SetConfigName(".rake")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("..")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func TestHighlight(t *testing.T) {
	instagram := shared.NewInstagram(viper.GetString("fbsr"), viper.GetString("session"), viper.GetString("app"))
	URLs, username, err := instagram.Reels("17898619759829276", true)
	assert.NoError(t, err)
	assert.Equal(t, "wikipedia", username)
	assert.Len(t, URLs, 8)
	for _, urlString := range URLs {
		URL, err := url.Parse(urlString)
		assert.NoError(t, err)
		assert.Equal(t, "https", URL.Scheme)
		assert.Regexp(t, instagramDomainRegularExpression, URL.Host, urlString)
		assert.Regexp(t, filePathRegularExpression, URL.Path)
	}
}

func TestStory(t *testing.T) {
	instagram := shared.NewInstagram(viper.GetString("fbsr"), viper.GetString("session"), viper.GetString("app"))
	URLs, username, err := instagram.Reels("f1", false)
	assert.NoError(t, err)
	assert.Equal(t, "f1", username)
	assert.Positive(t, len(URLs))
	for _, urlString := range URLs {
		URL, err := url.Parse(urlString)
		assert.NoError(t, err)
		assert.Equal(t, "https", URL.Scheme)
		assert.Regexp(t, instagramDomainRegularExpression, URL.Host, urlString)
		assert.Regexp(t, filePathRegularExpression, URL.Path)
	}
}
