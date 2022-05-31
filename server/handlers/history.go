package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/AppleGamer22/rake/server/db"
)

func History(writer http.ResponseWriter, request *http.Request) {
	webToken := request.Form.Get("token")
	if webToken == "" {
		http.Error(writer, "JWT must be provided", http.StatusBadRequest)
		return
	}

	payload, err := Authenticator.Parse(webToken)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		log.Println(err)
		return
	}

	result := db.Users.FindOne(context.Background(), db.User{ID: payload.U_ID})
	var user db.User
	if err := result.Decode(&user); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

}

func filterHistory(writer http.ResponseWriter, request *http.Request) {

}

func editHistory(writer http.ResponseWriter, request *http.Request) {

}

func deleteHistory() {

}

func HistoryPage(writer http.ResponseWriter, request *http.Request) {

}
