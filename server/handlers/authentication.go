package handlers

import (
	"context"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/AppleGamer22/rake/server/authenticator"
	"github.com/AppleGamer22/rake/server/db"
	"github.com/AppleGamer22/rake/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Authenticator authenticator.Authenticator

func InstagramSignUp(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		http.Error(writer, "failed to read request form", http.StatusBadRequest)
		return
	}

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

	secret := request.Form.Get("secret")
	if secret == "" {
		http.Error(writer, "secret must be provided", http.StatusBadRequest)
		return
	}

	count, err := db.Users.CountDocuments(context.Background(), bson.M{"username": username})
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	} else if count != 0 {
		http.Error(writer, "username already exists", http.StatusConflict)
		return
	}

	hashed, err := authenticator.Hash(password)
	if err != nil {
		http.Error(writer, "failed to store credentials securely", http.StatusInternalServerError)
		return
	}

	user := db.User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Hash:     hashed,
		Joined:   time.Now(),
		Network:  db.Instagram,
		// Instagram: false,
	}
	_, err = db.Users.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	webToken, err := Authenticator.Sign(authenticator.Payload{Username: user.Username, U_ID: user.ID})
	if err != nil {
		http.Error(writer, "sign-in failed", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	jwtCookie := &http.Cookie{
		Name:     "jwt",
		Value:    webToken,
		Path:     "/",
		Domain:   request.Host,
		HttpOnly: true,
	}

	http.SetCookie(writer, jwtCookie)

	instagramCookie := &http.Cookie{
		Name:     "instagram",
		Value:    password,
		Path:     "/",
		Domain:   request.Host,
		HttpOnly: true,
	}
	http.SetCookie(writer, instagramCookie)

}

func InstagramSignOut(writer http.ResponseWriter, request *http.Request) {
	jwtCookie, err := request.Cookie("jwt")
	if err != nil {
		http.Error(writer, "a JWT must be provided", http.StatusBadRequest)
		log.Println(err)
		return
	}

	instagramCookie, err := request.Cookie("instagram")
	if err != nil {
		http.Error(writer, "an Instagram cookie must be provided", http.StatusBadRequest)
		log.Println(err)
		return
	}

	payload, err := Authenticator.Parse(jwtCookie.Value)
	if err != nil {
		http.Error(writer, "sign-out failed", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	result := db.Users.FindOne(context.Background(), bson.M{"_id": payload.U_ID})
	var user db.User
	if err := result.Decode(&user); err != nil {
		http.Error(writer, "sign-out failed", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	if err := authenticator.Compare(user.Hash, instagramCookie.Value); err != nil {
		http.Error(writer, "sign-out failed", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	userDataDir := path.Join(shared.UserDataDirectory, user.ID.String())
	raker, err := shared.NewRaker(userDataDir, false, false)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if err := raker.InstagramSignOut(user.Username); err != nil {
		http.Error(writer, "sign-out failed", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if _, err := db.Users.DeleteOne(context.Background(), bson.M{"_id": user.ID}); err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		log.Println(err)
		return
	}

	jwtCookie.Value = ""
	jwtCookie.MaxAge = -1
	http.SetCookie(writer, jwtCookie)
	instagramCookie.Value = ""
	instagramCookie.MaxAge = -1
	http.SetCookie(writer, instagramCookie)
}

func AuthenticationPage(writer http.ResponseWriter, request *http.Request) {

}
