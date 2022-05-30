package authenticator

import (
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Authenticator struct {
	secret []byte
}

type Payload struct {
	jwt.RegisteredClaims
	Username string             `json:"username"`
	U_ID     primitive.ObjectID `json:"u_id"`
}

func New(secret string) Authenticator {
	return Authenticator{[]byte(secret)}
}

func (authenticator *Authenticator) Parse(tokenString string) (Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		return authenticator.secret, nil
	}, jwt.WithValidMethods(jwt.GetAlgorithms()))
	if payload, ok := token.Claims.(*Payload); ok && token.Valid {
		return *payload, nil
	} else {
		return Payload{}, err
	}
}

func (authenticator *Authenticator) Sign(payload Payload) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString(authenticator.secret)
}
