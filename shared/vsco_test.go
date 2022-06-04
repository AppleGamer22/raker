package shared_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/AppleGamer22/rake/shared"
	"github.com/stretchr/testify/assert"
)

func TestVSCOPicture(t *testing.T) {
	raker, err := shared.NewRaker("", "", false, true)
	assert.NoError(t, err)
	urlString, username, err := raker.VSCO("evgeneygolovesov", "6293acba42c9064de28f25b7")
	assert.NoError(t, err)
	assert.Equal(t, "evgeneygolovesov", username)
	URL, err := url.Parse(urlString)
	assert.NoError(t, err)
	assert.Equal(t, "https", URL.Scheme, urlString)
	assert.True(t, strings.Contains(URL.Host, "vsco.co"))
	assert.Regexp(t, filePathRegularExpression, URL.Path)
}
