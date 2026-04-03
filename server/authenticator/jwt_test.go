package authenticator_test

import (
	"testing"
	"time"

	"github.com/AppleGamer22/raker/server/authenticator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	testAuthenticator = authenticator.New("secret")
	expectedUsername  = uuid.NewString()
)

func TestJWT(t *testing.T) {
	webToken, expiry, err := testAuthenticator.Sign(expectedUsername)
	assert.NoError(t, err)
	assert.Less(t, expiry, time.Now().AddDate(1, 0, 0))
	username, err := testAuthenticator.Parse(webToken)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsername, username)
}
