package shared_test

import (
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

func testHighlight(t *testing.T) {
	instagram := shared.NewInstagram(viper.GetString("fbsr"), viper.GetString("session"), viper.GetString("app"))
	URLs, username, err := instagram.Reels("17898619759829276", true)
	assert.NoError(t, err)
	assert.Equal(t, "wikipedia", username)
	assert.Len(t, URLs, 8)
	testInstagramURLs(t, URLs)
}

func testStory(t *testing.T) {
	instagram := shared.NewInstagram(viper.GetString("fbsr"), viper.GetString("session"), viper.GetString("app"))
	URLs, username, err := instagram.Reels("f1", false)
	assert.NoError(t, err)
	assert.Equal(t, "f1", username)
	assert.Positive(t, len(URLs))
	testInstagramURLs(t, URLs)
}

func TestReels(t *testing.T) {
	t.Run("Highlight", testHighlight)
	t.Run("Story", testStory)
}
