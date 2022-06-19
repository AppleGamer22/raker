package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"

	"github.com/AppleGamer22/rake/server/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func History(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("jwt")
	if err != nil {
		http.Error(writer, "a JWT must be provided", http.StatusBadRequest)
		log.Println(err)
		return
	}

	payload, err := Authenticator.Parse(cookie.Value)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		log.Println(err)
		return
	}

	result := db.Users.FindOne(context.Background(), bson.M{"_id": payload.U_ID})
	var user db.User
	if err := result.Decode(&user); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := request.ParseForm(); err != nil {
		http.Error(writer, "failed to read request form", http.StatusBadRequest)
		return
	}
	media := request.Form.Get("media")
	owner := request.Form.Get("owner")
	post := request.Form.Get("post")
	file := request.URL.Query().Get("delete")
	categories := make([]string, 0, len(user.Categories))
	for _, category := range user.Categories {
		if request.Form.Get(category) == category {
			categories = append(categories, category)
		}
	}
	sort.Strings(categories)

	if media == "" || owner == "" || post == "" {
		http.Error(writer, "media type, owner & post must be valid", http.StatusBadRequest)
		return
	}

	switch request.Method {
	case http.MethodGet:
		histories, err := FilterHistories(user.ID, media, owner, post, categories)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			log.Println(err)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(histories)
	case http.MethodPost:
		switch request.Form.Get("method") {
		case http.MethodPatch:
			_, err := EditHistory(user.ID, media, owner, post, categories)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				log.Println(err)
				return
			}
		case http.MethodDelete:
			if file == "" {
				http.Error(writer, "file URL must be valid", http.StatusBadRequest)
				log.Println(err)
				return
			}

			_, err := DeleteFileFromHistory(user.ID, post, file)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				log.Println(err)
				return
			}
		default:
			http.Error(writer, "request method is not recognized", http.StatusBadRequest)
			return
		}
		http.Redirect(writer, request, request.Referer(), http.StatusTemporaryRedirect)
	default:
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func FilterHistories(U_ID primitive.ObjectID, media, owner, post string, categories []string) ([]db.History, error) {
	if !db.ValidMediaType(media) && media != "all" {
		return []db.History{}, errors.New("media must be valid")
	}

	filter := bson.M{"U_ID": U_ID.Hex()}

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
	err = cursor.All(context.Background(), &histories)
	return histories, err
}

func EditHistory(U_ID primitive.ObjectID, media, owner, post string, categories []string) (db.History, error) {
	filter := bson.M{
		"U_ID":  U_ID.Hex(),
		"type":  media,
		"owner": owner,
		"post":  post,
	}

	update := bson.M{
		"$set": bson.M{
			"categories": categories,
		},
	}

	var history db.History
	err := db.Histories.FindOneAndUpdate(context.Background(), filter, update).Decode(&history)
	return history, err
}

func DeleteFileFromHistory(U_ID primitive.ObjectID, post, file string) (db.History, error) {
	filter := bson.M{
		"U_ID": U_ID.Hex(),
		"urls": file,
		"post": post,
	}

	update := bson.M{
		"$pull": bson.M{
			"urls": file,
		},
	}

	result := db.Histories.FindOneAndUpdate(context.Background(), filter, update)
	var history db.History
	if err := result.Decode(&history); err != nil {
		return db.History{}, err
	}

	if len(history.URLs) == 0 {
		result, err := db.Histories.DeleteOne(context.Background(), filter)
		if err != nil {
			return db.History{}, err
		} else if result.DeletedCount == 0 {
			return db.History{}, errors.New("no histories were found")
		} else {
			return db.History{}, nil
		}
	}

	return history, nil
}

func HistoryPage(writer http.ResponseWriter, request *http.Request) {

}
