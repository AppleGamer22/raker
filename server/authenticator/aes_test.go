package authenticator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const expectedPlainText = "an important secret"

func TestAES(t *testing.T) {
	ct, err := testAuthenticator.Encrypt(expectedPlainText)
	assert.NoError(t, err)
	pt, err := testAuthenticator.Decrypt(ct)
	assert.NoError(t, err)
	assert.Equal(t, expectedPlainText, pt)
}
