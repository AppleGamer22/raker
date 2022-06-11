package authenticator_test

import (
	"testing"

	"github.com/AppleGamer22/rake/server/authenticator"
	"github.com/stretchr/testify/assert"
)

const (
	expectedPlainText = "a secret"
	key               = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOiJhIGtleSJ9.pF8Mr1K7FECoiWtCjO-IIez2s94iwFAaF3VRyljXloU"
)

func TestAES(t *testing.T) {
	ct, err := authenticator.Encrypt(key, expectedPlainText)
	assert.NoError(t, err)
	pt, err := authenticator.Decrypt(key, ct)
	assert.NoError(t, err)
	assert.Equal(t, expectedPlainText, pt)
}
