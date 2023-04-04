package handlers

import (
	"context"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/AppleGamer22/raker/server/authenticator"
	"github.com/AppleGamer22/raker/server/cleaner"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/AppleGamer22/raker/shared"
	"github.com/AppleGamer22/raker/shared/types"
	"github.com/charmbracelet/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Authenticator authenticator.Authenticator

func InstagramSignUp(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		return
	}

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
	if fbsr == "" {
		http.Error(writer, "FBSR must be provided", http.StatusBadRequest)
		return
	}

	sessionID := cleaner.Line(request.Form.Get("session"))
	if sessionID == "" {
		http.Error(writer, "session ID must be provided", http.StatusBadRequest)
		return
	}

	userID := cleaner.Line(request.Form.Get("user"))
	if userID == "" {
		http.Error(writer, "user ID must be provided", http.StatusBadRequest)
		return
	}

	count, err := db.Users.CountDocuments(context.Background(), bson.M{"username": username})
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Error(err, username)
		return
	} else if count != 0 {
		http.Error(writer, "username already exists", http.StatusConflict)
		log.Error("username already exists", username)
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
		Network: types.Instagram,
	}
	_, err = db.Users.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	InstagramSignIn(writer, request)
}

func InstagramSignIn(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		return
	}

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
		log.Error(err)
		return
	}

	if err := authenticator.Compare(user.Hash, password); err != nil {
		http.Error(writer, "incorrect credentials", http.StatusUnauthorized)
		log.Error(err)
		return
	}

	webToken, expiry, err := Authenticator.Sign(user.ID, user.Username)
	if err != nil {
		http.Error(writer, "sign-in failed", http.StatusUnauthorized)
		log.Error(err, user.ID.Hex())
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:     "jwt",
		Value:    webToken,
		Path:     "/",
		Domain:   request.Host,
		Expires:  expiry,
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

	U_ID, username, err := Authenticator.Parse(jwtCookie.Value)
	if err != nil {
		return db.User{}, err
	}

	filter := bson.M{
		"_id":      U_ID,
		"username": username,
	}
	var user db.User
	err = db.Users.FindOne(context.Background(), filter).Decode(&user)
	return user, err
}

func InstagramUpdateCredentials(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		return
	}

	user, err := Verify(request)
	if err != nil {
		http.Error(writer, "credential update failed", http.StatusUnauthorized)
		log.Error(err, user.ID.Hex())
		return
	}

	if err := request.ParseForm(); err != nil {
		http.Error(writer, "failed to read request form", http.StatusBadRequest)
		return
	}

	password := cleaner.Line(request.Form.Get("password"))
	if password == "" {
		password = user.Hash
	} else {
		password, err = authenticator.Hash(password)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			log.Error(err, user.ID.Hex())
			return
		}
	}

	fbsr := cleaner.Line(request.Form.Get("fbsr"))
	if fbsr == "" {
		fbsr = user.Instagram.FBSR
	}

	sessionID := cleaner.Line(request.Form.Get("session"))
	if sessionID == "" {
		sessionID = user.Instagram.SessionID
	}

	userID := cleaner.Line(request.Form.Get("user"))
	if userID == "" {
		userID = user.Instagram.UserID
	}

	filter := bson.M{
		"_id":      user.ID,
		"username": user.Username,
	}

	update := bson.M{
		"$set": bson.M{
			"hash":                 password,
			"instagram.fbsr":       fbsr,
			"instagram.session_id": sessionID,
			"instagram.user_id":    userID,
		},
	}

	result, err := db.Users.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Error(err, user.ID.Hex())
		return
	} else if result.MatchedCount == 0 {
		http.Error(writer, "the user was not found/modified", http.StatusNotFound)
		log.Error("the user was not found/modified", user.ID.Hex())
	}

	http.Redirect(writer, request, request.Referer(), http.StatusTemporaryRedirect)
}

func InstagramSignOut(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		return
	}

	user, err := Verify(request)
	if err != nil {
		http.Error(writer, "sign-out failed", http.StatusUnauthorized)
		log.Error(err, user.ID.Hex())
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:   "jwt",
		Value:  "",
		MaxAge: -1,
	})

	http.Redirect(writer, request, request.Referer(), http.StatusTemporaryRedirect)
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

	if err := templates.ExecuteTemplate(writer, "authentication.html", categoryDisplay); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
}
