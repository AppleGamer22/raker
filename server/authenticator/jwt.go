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
	Username string `json:"username"`
	// Password string             `json:"password"`
	U_ID primitive.ObjectID `json:"U_ID"`
}

func New(secret string) Authenticator {
	return Authenticator{[]byte(secret)}
}

func (a *Authenticator) Parse(tokenString string) (Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		return a.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if payload, ok := token.Claims.(*Payload); ok && token.Valid {
		return *payload, nil
	} else {
		return Payload{}, err
	}
}

func (a *Authenticator) Sign(payload Payload) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString(a.secret)
}
