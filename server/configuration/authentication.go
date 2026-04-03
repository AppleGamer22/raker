package configuration

import (
	"context"
	"html/template"
	"net/http"
	"net/url"

	"github.com/AppleGamer22/raker/server/authenticator"
	"github.com/AppleGamer22/raker/server/cleaner"
	"github.com/AppleGamer22/raker/server/db"
	old "github.com/AppleGamer22/raker/server/db/mongo"

	"github.com/AppleGamer22/raker/shared"
	"github.com/AppleGamer22/raker/templates"
	"github.com/charmbracelet/log"
)

func (server *RakerServer) InstagramSignUp(writer http.ResponseWriter, request *http.Request) {
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

	_, err := server.DBClient.UserGet(context.Background(), username)
	if err == nil {
		http.Error(writer, "username already exists", http.StatusConflict)
		log.Error("username already exists", "username", username)
		return
	}

	hashed, err := authenticator.Hash(password)
	if err != nil {
		http.Error(writer, "failed to store credentials securely", http.StatusInternalServerError)
		return
	}

	err = server.DBClient.UserAdd(context.Background(), db.UserAddParams{
		Username:           username,
		PasswordHash:       hashed,
		InstagramSessionID: sessionID,
		InstagramUserID:    userID,
	})
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}

	server.InstagramSignIn(writer, request)
}

func (server *RakerServer) InstagramSignIn(writer http.ResponseWriter, request *http.Request) {
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

	user, err := server.DBClient.UserGet(context.Background(), username)
	if err != nil {
		http.Error(writer, "incorrect credentials", http.StatusUnauthorized)
		log.Error(err)
		return
	}

	if err := authenticator.Compare(user.PasswordHash, password); err != nil {
		http.Error(writer, "incorrect credentials", http.StatusUnauthorized)
		log.Error(err)
		return
	}

	webToken, expiry, err := server.Authenticator.Sign(user.Username)
	if err != nil {
		http.Error(writer, "sign-in failed", http.StatusUnauthorized)
		log.Error(err, "ID", user.Username)
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:  "jwt",
		Value: webToken,
		Path:  "/",
		// Domain:   request.URL.Hostname(),
		Expires:  expiry,
		Secure:   server.Configuration.SecureCookie,
		HttpOnly: true,
	})

	http.Redirect(writer, request, request.Referer(), http.StatusTemporaryRedirect)
}

func (server *RakerServer) getUserFromCookie(request *http.Request) (db.User, error) {
	jwtCookie, err := request.Cookie("jwt")
	if err != nil {
		return db.User{}, err
	}

	username, err := server.Authenticator.Parse(jwtCookie.Value)
	if err != nil {
		return db.User{}, err
	}

	user, err := server.DBClient.UserGet(context.Background(), username)
	return user, err
}

func (server *RakerServer) Verify(strict bool, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user, err := server.getUserFromCookie(request)
		if err != nil {
			log.Error(err)
			if strict {
				http.Error(writer, "credential verification failed", http.StatusUnauthorized)
				return
			}
		}
		// https://drstearns.github.io/tutorials/gomiddleware/#secmiddlewareandrequestscopedvalues
		ctxWithUser := context.WithValue(request.Context(), authenticatedUserKey, user)
		requestWithUser := request.WithContext(ctxWithUser)
		handler.ServeHTTP(writer, requestWithUser)
	})
}

func (server *RakerServer) InstagramUpdateCredentials(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		return
	}

	user := request.Context().Value(authenticatedUserKey).(db.User)

	err := request.ParseForm()
	if err != nil {
		http.Error(writer, "failed to read request form", http.StatusBadRequest)
		return
	}

	password := cleaner.Line(request.Form.Get("password"))
	if password == "" {
		password = user.PasswordHash
	} else {
		password, err = authenticator.Hash(password)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			log.Error(err, "ID", user.Username)
			return
		}
	}

	// fbsr := cleaner.Line(request.Form.Get("fbsr"))
	// if fbsr == "" {
	// 	fbsr = user.Instagram.FBSR
	// }

	sessionID := cleaner.Line(request.Form.Get("session"))
	if sessionID == "" {
		sessionID = user.InstagramSessionID
	}

	userID := cleaner.Line(request.Form.Get("user"))
	if userID == "" {
		userID = user.InstagramUserID
	}

	err = server.DBClient.UserUpdateInstagramSession(context.Background(), db.UserUpdateInstagramSessionParams{
		InstagramSessionID: sessionID,
		InstagramUserID:    userID,
		Username:           user.Username,
	})
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Error(err, "ID", user.Username)
		return
	}

	http.Redirect(writer, request, request.Referer(), http.StatusTemporaryRedirect)
}

func (server *RakerServer) InstagramSignOut(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		return
	}

	// user := request.Context().Value(authenticatedUserKey).(db.User)

	http.SetCookie(writer, &http.Cookie{
		Name:   "jwt",
		Value:  "",
		MaxAge: -1,
	})

	http.Redirect(writer, request, request.Referer(), http.StatusTemporaryRedirect)
}

func (server *RakerServer) AuthenticationPage(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value(authenticatedUserKey).(db.User)

	categoryDisplay := old.UserCategoryDisplay{
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

	if err := templates.Templates.ExecuteTemplate(writer, "authentication.html", categoryDisplay); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
}
