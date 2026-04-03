package authenticator

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type Authenticator struct {
	secret []byte
}

type jwtPayload struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}

func New(secret string) Authenticator {
	return Authenticator{[]byte(secret)}
}

func (a *Authenticator) Parse(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtPayload{}, func(token *jwt.Token) (interface{}, error) {
		return a.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if payload, ok := token.Claims.(*jwtPayload); ok && token.Valid {
		return payload.Username, nil
	} else {
		return "", err
	}
}

func (a *Authenticator) Sign(username string) (string, time.Time, error) {
	payload := jwtPayload{
		Username: username,
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
