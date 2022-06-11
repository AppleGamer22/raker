package authenticator_test

import (
	"testing"

	"github.com/AppleGamer22/rake/server/authenticator"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var testAuthenticator = authenticator.New("secret")

var expectedPayload = authenticator.Payload{
	Username: "rake",
	U_ID: func() primitive.ObjectID {
		_id, _ := primitive.ObjectIDFromHex("6294257e01948be67fe71b94")
		return _id
	}(),
}

func TestJWT(t *testing.T) {
	webToken, err := testAuthenticator.Sign(expectedPayload)
	assert.NoError(t, err)
	payload, err := testAuthenticator.Parse(webToken)
	assert.NoError(t, err)
	assert.Equal(t, expectedPayload.Username, payload.Username)
	assert.Equal(t, expectedPayload.U_ID, payload.U_ID)
	assert.Equal(t, expectedPayload, payload)
}
