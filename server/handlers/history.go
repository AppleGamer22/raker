package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/AppleGamer22/rake/server/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	media := request.Form.Get("media")
	owner := request.Form.Get("owner")
	post := request.Form.Get("post")
	categories := strings.Split(request.Form.Get("categories"), ",")

	if media == "" || owner == "" || post == "" {
		http.Error(writer, "media type, owner & post must be valid", http.StatusBadRequest)
		return
	}

	switch request.Method {
	case http.MethodGet:
		histories, err := filterHistory(user.ID, media, owner, post, categories)
		if err != nil {
			http.Error(writer, "media must be valid", http.StatusBadRequest)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(histories)
	case http.MethodPatch:
	default:
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func filterHistory(U_ID primitive.ObjectID, media, owner, post string, categories []string) ([]db.History, error) {
	if !db.ValidMediaType(media) && media != "all" {
		return []db.History{}, errors.New("media must be valid")
	}

	filter := bson.M{"U_ID": U_ID.String()}

	if len(categories) != 0 {
		filter["categories"] = categories
	} else {
		filter["categories"] = bson.M{"$in": primitive.A{primitive.Undefined{}, primitive.A{}}}
	}

	if owner != "all" {
		filter["owner"] = bson.M{"$regex": primitive.Regex{
			Pattern: owner,
			Options: "i",
		}}
	}

	cursor, err := db.Histories.Find(context.Background(), filter)
	if err != nil {
		return []db.History{}, err
	}

	var histories []db.History
	if err := cursor.All(context.Background(), &histories); err != nil {
		return []db.History{}, err
	}
	return histories, nil
}

func editHistory(writer http.ResponseWriter, request *http.Request) {

}

func deleteHistory() {

}

func HistoryPage(writer http.ResponseWriter, request *http.Request) {

}
