package shared_test

import (
	"net/url"
	"testing"

	"github.com/AppleGamer22/raker/shared"
	"github.com/stretchr/testify/assert"
)

func TestSnapchatHighlight(t *testing.T) {
	result, _, err := shared.SnapchatHighlight("nasa", "264684ee-96ec-4281-a1a9-79562a350f4e")
	assert.NoError(t, err)
	assert.Equal(t, "nasa", result.Username)
	assert.NotEmpty(t, result.URLs)
	assert.False(t, result.IsStory)

	for _, item := range result.URLs {
		parsedURL, err := url.Parse(item.URL)
		assert.NoError(t, err)
		assert.Equal(t, "https", parsedURL.Scheme)
		assert.Contains(t, parsedURL.Host, "sc-cdn.net")
	}
}
