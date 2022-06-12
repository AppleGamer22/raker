package handlers

import (
	"context"
	"net/http"

	"github.com/AppleGamer22/rake/server/db"
	"go.mongodb.org/mongo-driver/bson"
)

func Verify(request *http.Request) (db.User, error) {
	jwtCookie, err := request.Cookie("jwt")
	if err != nil {
		return db.User{}, err
	}

	payload, err := Authenticator.Parse(jwtCookie.Value)
	if err != nil {
		return db.User{}, err
	}

	var user db.User
	err = db.Users.FindOne(context.Background(), bson.M{"_id": payload.U_ID}).Decode(&user)
	return user, err
}
