package authenticator

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Authenticator struct {
	secret []byte
}

type jwtPayload struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	// Password string             `json:"password"`
	U_ID primitive.ObjectID `json:"U_ID"`
}

func New(secret string) Authenticator {
	return Authenticator{[]byte(secret)}
}

func (a *Authenticator) Parse(tokenString string) (primitive.ObjectID, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtPayload{}, func(token *jwt.Token) (interface{}, error) {
		return a.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if payload, ok := token.Claims.(*jwtPayload); ok && token.Valid {
		return payload.U_ID, payload.Username, nil
	} else {
		return primitive.NilObjectID, "", err
	}
}

func (a *Authenticator) Sign(U_ID primitive.ObjectID, username string) (string, time.Time, error) {
	payload := jwtPayload{
		Username: username,
		U_ID:     U_ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(1, 0, 0)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	jsonWebToken, err := token.SignedString(a.secret)
	if err != nil {
		return "", time.Unix(0, 0), err
	}
	return jsonWebToken, payload.ExpiresAt.Time, nil
}
