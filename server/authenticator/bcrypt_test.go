package authenticator_test

import (
	"testing"

	"github.com/AppleGamer22/raker/server/authenticator"
	"github.com/stretchr/testify/assert"
)

const (
	expectedPassword = "@!X*WHW6@^yRa!JZh#REMu28Mv^ae"
	expectedHash     = "$2a$10$vGYdgqsg6MujCSNeGOTu/uQwalug9MFbIvCl0f8v1WOulD99phRIy"
)

func TestBcrypt(t *testing.T) {
	t.Run("Hash", func(t *testing.T) {
		hashed, err := authenticator.Hash(expectedPassword)
		assert.NoError(t, err)
		assert.NoError(t, authenticator.Compare(hashed, expectedPassword))
	})

	t.Run("Compare", func(t *testing.T) {
		assert.NoError(t, authenticator.Compare(expectedHash, expectedPassword))
	})
}
