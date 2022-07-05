package handlers

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/AppleGamer22/rake/server/authenticator"
	"github.com/AppleGamer22/rake/server/cleaner"
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

	username := cleaner.Line(request.Form.Get("username"))
	if username == "" {
		http.Error(writer, "username must be provided", http.StatusBadRequest)
		return
	}

	password := cleaner.Line(request.Form.Get("password"))
	if password == "" {
		http.Error(writer, "password must be provided", http.StatusBadRequest)
		return
	}

	fbsr := cleaner.Line(request.Form.Get("fbsr"))
	if password == "" {
		http.Error(writer, "FBSR must be provided", http.StatusBadRequest)
		return
	}

	sessionID := cleaner.Line(request.Form.Get("session"))
	if password == "" {
		http.Error(writer, "session ID must be provided", http.StatusBadRequest)
		return
	}

	userID := cleaner.Line(request.Form.Get("user"))
	if password == "" {
		http.Error(writer, "user ID must be provided", http.StatusBadRequest)
		return
	}

	count, err := db.Users.CountDocuments(context.Background(), bson.M{"username": username})
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Println(err)
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
		Instagram: struct {
			FBSR      string `bson:"fbsr" json:"-"`
			SessionID string `bson:"session_id" json:"-"`
			UserID    string `bson:"user_id" json:"-"`
			// AppID     string `bson:"app_id" json:"-"`
		}{
			FBSR:      fbsr,
			SessionID: sessionID,
			UserID:    userID,
		},
		Joined:  time.Now(),
		Network: db.Instagram,
	}
	_, err = db.Users.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Println()
		return
	}

	InstagramSignIn(writer, request)
}

func InstagramSignIn(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		http.Error(writer, "failed to read request form", http.StatusBadRequest)
		return
	}

	username := cleaner.Line(request.Form.Get("username"))
	if username == "" {
		http.Error(writer, "username must be provided", http.StatusBadRequest)
		return
	}

	password := cleaner.Line(request.Form.Get("password"))
	if password == "" {
		http.Error(writer, "password must be provided", http.StatusBadRequest)
		return
	}

	var user db.User
	if err := db.Users.FindOne(context.Background(), bson.M{"username": username}).Decode(&user); err != nil {
		http.Error(writer, "incorrect credentials", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	if err := authenticator.Compare(user.Hash, password); err != nil {
		http.Error(writer, "incorrect credentials", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	webToken, err := Authenticator.Sign(authenticator.Payload{Username: user.Username, U_ID: user.ID})
	if err != nil {
		http.Error(writer, "sign-in failed", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:     "jwt",
		Value:    webToken,
		Path:     "/",
		Domain:   request.Host,
		Expires:  time.Now().AddDate(1, 0, 0),
		Secure:   true,
		HttpOnly: true,
	})

	http.Redirect(writer, request, request.Referer(), http.StatusTemporaryRedirect)
}

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

func InstagramSignOut(writer http.ResponseWriter, request *http.Request) {
	_, err := Verify(request)
	if err != nil {
		http.Error(writer, "sign-out failed", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:   "jwt",
		Value:  "",
		MaxAge: -1,
	})
}

func AuthenticationPage(writer http.ResponseWriter, request *http.Request) {
	user, _ := Verify(request)

	categoryDisplay := db.UserCategoryDisplay{
		Username:   user.Username,
		Categories: user.Categories,
		HistoryQuery: func() template.URL {
			query := url.Values{}
			query.Set("page", "1")
			for _, category := range user.Categories {
				query.Set(category, category)
			}
			return template.URL(query.Encode())
		}(),
		Version: shared.Version,
	}

	writer.Header().Set("Content-Type", "text/html")
	if err := templates.ExecuteTemplate(writer, "authentication.html", categoryDisplay); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
