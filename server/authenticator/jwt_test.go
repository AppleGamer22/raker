package authenticator_test

import (
	"testing"
	"time"

	"github.com/AppleGamer22/rake/server/authenticator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	testAuthenticator = authenticator.New("secret")
	expectedUsername  = uuid.NewString()
	expectedUserID    = primitive.NewObjectID()
)

func TestJWT(t *testing.T) {
	webToken, expiry, err := testAuthenticator.Sign(expectedUserID, expectedUsername)
	assert.NoError(t, err)
	assert.Less(t, expiry, time.Now().AddDate(1, 0, 0))
	U_ID, username, err := testAuthenticator.Parse(webToken)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsername, username)
	assert.Equal(t, expectedUserID, U_ID)
}
