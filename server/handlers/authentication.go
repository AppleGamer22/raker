package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AppleGamer22/rake/server/authenticator"
	"github.com/AppleGamer22/rake/server/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Authenticator authenticator.Authenticator

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	WebToken string `json:"token"`
}

func InstagramSignUp(writer http.ResponseWriter, request *http.Request) {
	username := request.Form.Get("username")
	if username == "" {
		http.Error(writer, "username must be provided", http.StatusBadRequest)
		return
	}

	password := request.Form.Get("password")
	if password == "" {
		http.Error(writer, "password must be provided", http.StatusBadRequest)
		return
	}

	count, err := db.Users.CountDocuments(context.Background(), db.User{Username: username})
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	} else if count != 0 {
		http.Error(writer, "username already exists", http.StatusConflict)
		return
	}

	hashed, err := authenticator.Hash(password)
	if err != nil {
		http.Error(writer, "failed to store password securely", http.StatusInternalServerError)
		return
	}
	user := db.User{
		ID:        primitive.NewObjectID(),
		Username:  username,
		Hash:      hashed,
		Joined:    time.Now(),
		Network:   db.Instagram,
		Instagram: false,
	}
	result, err := db.Users.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(writer, result.InsertedID)

}

func InstagramSignIn(writer http.ResponseWriter, request *http.Request) {
	username := request.Form.Get("username")
	if username == "" {
		http.Error(writer, "username must be provided", http.StatusBadRequest)
		return
	}

	password := request.Form.Get("password")
	if password == "" {
		http.Error(writer, "password must be provided", http.StatusBadRequest)
		return
	}

	result := db.Users.FindOne(context.Background(), db.User{Username: username})
	var user db.User
	if err := result.Decode(&user); err != nil {
		http.Error(writer, "sign-in failed", http.StatusBadRequest)
		log.Println(err)
		return
	}

	if err := authenticator.Compare(user.Hash, password); err != nil {
		http.Error(writer, "sign-in failed", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	user.Instagram = true
	if _, err := db.Users.UpdateOne(context.Background(), db.User{ID: user.ID}, user); err != nil {
		http.Error(writer, "sign-in failed", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	webToken, err := Authenticator.Sign(authenticator.Payload{Username: user.Username, U_ID: user.ID})
	if err != nil {
		http.Error(writer, "sign-in failed", http.StatusUnauthorized)
		log.Println(err)
		return
	}
	// TODO: sign-in instagram
	fmt.Fprint(writer, webToken)
}

func InstagramSignOut(writer http.ResponseWriter, request *http.Request) {
	webToken := request.Form.Get("token")
	if webToken == "" {
		http.Error(writer, "JWT must be provided", http.StatusBadRequest)
		return
	}

	payload, err := Authenticator.Parse(webToken)
	if err != nil {
		http.Error(writer, "sign-out failed", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	result := db.Users.FindOne(context.Background(), db.User{ID: payload.U_ID})
	var user db.User
	if err := result.Decode(&user); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: sign-out instagram
	user.Instagram = false
	if _, err := db.Users.UpdateOne(context.Background(), db.User{ID: user.ID}, user); err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		log.Println(err)
		return
	}
	fmt.Fprint(writer, "signed-out")
}

func AuthenticationPage(writer http.ResponseWriter, request *http.Request) {

}
