package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/AppleGamer22/rake/server/authenticator"
	"github.com/AppleGamer22/rake/server/db"
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

	session := request.Form.Get("session")
	if session == "" {
		http.Error(writer, "password must be provided", http.StatusBadRequest)
		return
	}

	fbsr := request.Form.Get("fbsr")
	if fbsr == "" {
		http.Error(writer, "an FBSR must be provided", http.StatusBadRequest)
		return
	}

	appID := request.Form.Get("app")
	if appID == "" {
		http.Error(writer, "an app ID must be provided", http.StatusBadRequest)
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

	hashedSession, err := authenticator.Hash(session)
	if err != nil {
		http.Error(writer, "failed to store credentials securely", http.StatusInternalServerError)
		return
	}

	hashedFBSR, err := authenticator.Hash(fbsr)
	if err != nil {
		http.Error(writer, "failed to store credentials securely", http.StatusInternalServerError)
		return
	}

	hashedAppID, err := authenticator.Hash(appID)
	if err != nil {
		http.Error(writer, "failed to store credentials securely", http.StatusInternalServerError)
		return
	}

	user := db.User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Session:  hashedSession,
		FBSR:     hashedFBSR,
		AppID:    hashedAppID,
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

	sessionCookie := &http.Cookie{
		Name:     "session",
		Value:    session,
		Path:     "/",
		Domain:   request.Host,
		HttpOnly: true,
	}
	http.SetCookie(writer, sessionCookie)

	fbsrCookie := &http.Cookie{
		Name:     "fbsr",
		Value:    fbsr,
		Path:     "/",
		Domain:   request.Host,
		HttpOnly: true,
	}
	http.SetCookie(writer, fbsrCookie)

	appIDCookie := &http.Cookie{
		Name:     "app_id",
		Value:    appID,
		Path:     "/",
		Domain:   request.Host,
		HttpOnly: true,
	}
	http.SetCookie(writer, appIDCookie)

}

func InstagramSignOut(writer http.ResponseWriter, request *http.Request) {
	user, err := Verify(request)
	if err != nil {
		http.Error(writer, "sign-out failed", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	if _, err := db.Users.DeleteOne(context.Background(), bson.M{"_id": user.ID}); err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		log.Println(err)
		return
	}

	jwtCookie := &http.Cookie{
		Name:   "jwt",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(writer, jwtCookie)

	sessionCookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(writer, sessionCookie)

	fbsrCookie := &http.Cookie{
		Name:   "fbsr",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(writer, fbsrCookie)

	appIDCookie := &http.Cookie{
		Name:   "app_id",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(writer, appIDCookie)
}

func AuthenticationPage(writer http.ResponseWriter, request *http.Request) {

}
