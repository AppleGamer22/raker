package shared_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/AppleGamer22/raker/shared"
	"github.com/stretchr/testify/assert"
)

func TestVSCOPicture(t *testing.T) {
	urlStrings, username, _, err := shared.VSCO("evgeneygolovesov", "6293acba42c9064de28f25b7")
	assert.NoError(t, err)
	assert.Equal(t, "evgeneygolovesov", username)
	assert.Len(t, urlStrings, 1)
	urlString := urlStrings[0]
	URL, err := url.Parse(urlString)
	assert.NoError(t, err)
	assert.Equal(t, "https", URL.Scheme, urlString)
	assert.True(t, strings.Contains(URL.Host, "vsco.co"), urlString)
	assert.Regexp(t, filePathRegularExpression, URL.Path)
}
