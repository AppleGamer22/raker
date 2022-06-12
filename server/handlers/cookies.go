package handlers

import (
	"context"
	"net/http"

	"github.com/AppleGamer22/rake/server/authenticator"
	"github.com/AppleGamer22/rake/server/db"
	"go.mongodb.org/mongo-driver/bson"
)

func Verify(request *http.Request) (db.User, error) {
	jwtCookie, err := request.Cookie("jwt")
	if err != nil {
		return db.User{}, err
	}

	sessionCookie, err := request.Cookie("session")
	if err != nil {
		return db.User{}, err
	}

	fbsrCookie, err := request.Cookie("fbsr")
	if err != nil {
		return db.User{}, err
	}

	appIDCookie, err := request.Cookie("app_id")
	if err != nil {
		return db.User{}, err
	}

	payload, err := Authenticator.Parse(jwtCookie.Value)
	if err != nil {
		return db.User{}, err
	}

	result := db.Users.FindOne(context.Background(), bson.M{"_id": payload.U_ID})
	var user db.User
	if err := result.Decode(&user); err != nil {
		return db.User{}, err
	}

	if err := authenticator.Compare(user.Session, sessionCookie.Value); err != nil {
		return db.User{}, err
	}

	if err := authenticator.Compare(user.FBSR, fbsrCookie.Value); err != nil {
		return db.User{}, err
	}

	if err := authenticator.Compare(user.AppID, appIDCookie.Value); err != nil {
		return db.User{}, err
	}

	return user, nil
}
