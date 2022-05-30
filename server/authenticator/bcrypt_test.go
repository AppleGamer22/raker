package authenticator_test

import (
	"testing"

	"github.com/AppleGamer22/rake/server/authenticator"
	"github.com/stretchr/testify/assert"
)

const (
	expectedPassword = "@!X*WHW6@^yRa!JZh#REMu28Mv^ae"
	expectedHash     = "$2a$10$vGYdgqsg6MujCSNeGOTu/uQwalug9MFbIvCl0f8v1WOulD99phRIy"
)

func TestHash(t *testing.T) {
	hashed, err := authenticator.Hash(expectedPassword)
	assert.NoError(t, err)
	assert.NoError(t, authenticator.Compare(hashed, expectedPassword))
}

func TestCompare(t *testing.T) {
	assert.NoError(t, authenticator.Compare(expectedHash, expectedPassword))
}
