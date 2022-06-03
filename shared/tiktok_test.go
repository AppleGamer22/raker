package shared_test

import (
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
	err := viper.ReadInConfig()
	assert.NoError(t, err)
	raker, err := shared.NewRaker(shared.FindExecutablePath(), "", false, true)
	assert.NoError(t, err)
	URL, username, err := raker.TikTok("f1", "7048983181063687430")
	assert.NoError(t, err)
	assert.Equal(t, "f1", username)
	assert.True(t, strings.HasPrefix(URL, "https://"), URL)
	assert.True(t, strings.Contains(URL, "-webapp.tiktok.com"), URL)
}
