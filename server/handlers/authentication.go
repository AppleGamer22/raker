package handlers

import (
	"context"
	"encoding/json"
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
	var creds credentials
	if err := json.NewDecoder(request.Body).Decode(&creds); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	count, err := db.Users.CountDocuments(context.TODO(), db.User{Username: creds.Username})
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	} else if count != 0 {
		http.Error(writer, "username already exists", http.StatusConflict)
		return
	}

	hashed, err := authenticator.Hash(creds.Password)
	if err != nil {
		http.Error(writer, "failed to store password securely", http.StatusInternalServerError)
		return
	}
	user := db.User{
		ID:        primitive.NewObjectID(),
		Hash:      hashed,
		Username:  creds.Username,
		Joined:    time.Now(),
		Network:   db.Instagram,
		Instagram: false,
	}
	result, err := db.Users.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(writer, result.InsertedID)

}

func InstagramSignIn(writer http.ResponseWriter, request *http.Request) {
	var creds credentials
	if err := json.NewDecoder(request.Body).Decode(&creds); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}
	result := db.Users.FindOne(context.TODO(), db.User{Username: creds.Username})
	var user db.User
	if err := result.Decode(&user); err != nil {
		http.Error(writer, "sign-in failed", http.StatusBadRequest)
		log.Println(err)
		return
	}

	if err := authenticator.Compare(user.Hash, creds.Password); err != nil {
		http.Error(writer, "sign-in failed", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	user.Instagram = true
	if _, err := db.Users.UpdateOne(context.TODO(), db.User{ID: user.ID}, user); err != nil {
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
	var creds credentials
	if err := json.NewDecoder(request.Body).Decode(&creds); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	payload, err := Authenticator.Parse(creds.WebToken)
	if err != nil {
		http.Error(writer, "sign-out failed", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	count, err := db.Users.CountDocuments(context.TODO(), db.User{ID: payload.U_ID})
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	} else if count == 0 {
		http.Error(writer, "username does not exist", http.StatusUnauthorized)
		return
	}
	// TODO: sign-out instagram
	fmt.Fprint(writer, "signed-out")
}

func AuthenticationPage(writer http.ResponseWriter, request *http.Request) {

}
