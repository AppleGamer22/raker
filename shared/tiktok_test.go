package shared_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/AppleGamer22/rake/shared"
	"github.com/stretchr/testify/assert"
)

func TestTikTokPublicVideo(t *testing.T) {
	tiktok := shared.NewTikTok("")
	urlString, username, err := tiktok.Post("f1", "7048983181063687430", false)
	assert.NoError(t, err)
	assert.Equal(t, "f1", username)
	URL, err := url.Parse(urlString)
	assert.NoError(t, err)
	assert.Equal(t, "https", URL.Scheme)
	assert.True(t, strings.Contains(URL.Host, "-webapp.tiktok.com"), urlString)
}
