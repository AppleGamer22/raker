package authenticator_test

import (
	"testing"

	"github.com/AppleGamer22/rake/server/authenticator"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var testAuthenticator = authenticator.New("secret")

const expectedToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IkFwcGxlR2FtZXIyMiIsInVfaWQiOiI2Mjk0MjU3ZTAxOTQ4YmU2N2ZlNzFiOTQifQ.rNPA79EHMMikHQbF5Tw1HYqLRkN-bIqTRwZnBqUq5Xs"

var expectedPayload = authenticator.Payload{
	Username: "AppleGamer22",
	U_ID: func() primitive.ObjectID {
		_id, _ := primitive.ObjectIDFromHex("6294257e01948be67fe71b94")
		return _id
	}(),
}

func TestSign(t *testing.T) {
	token, err := testAuthenticator.Sign(expectedPayload)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
}

func TestParse(t *testing.T) {
	payload, err := testAuthenticator.Parse(expectedToken)
	assert.NoError(t, err)
	assert.Equal(t, expectedPayload.Username, payload.Username)
	assert.Equal(t, expectedPayload.U_ID, payload.U_ID)
	assert.Equal(t, expectedPayload, payload)
}
