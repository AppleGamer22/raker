package authenticator

import (
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Authenticator struct {
	secret []byte
}

type Payload struct {
	Username string
	ID       primitive.ObjectID
}

func New(secret string) Authenticator {
	return Authenticator{[]byte(secret)}
}

func (authenticator *Authenticator) Parse(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return authenticator.secret, nil
	}, jwt.WithValidMethods(jwt.GetAlgorithms()))
}

func (authenticator *Authenticator) Sign(payload Payload) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Username": payload.Username,
		"ID":       payload.ID.String(),
	})
	return token.SignedString(authenticator.secret)
}
