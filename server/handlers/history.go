package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"

	"github.com/AppleGamer22/rake/server/cleaner"
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
	media := cleaner.Line(request.Form.Get("media"))
	owner := cleaner.Line(request.Form.Get("owner"))
	post := cleaner.Line(request.Form.Get("post"))
	file := cleaner.Line(request.Form.Get("remove"))

	categories := make([]string, 0, len(user.Categories))
	for _, category := range user.Categories {
		if cleaner.Line(request.Form.Get(category)) == category {
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
		histories, err := filterHistories(user.ID, media, owner, post, categories)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			log.Println(err)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(histories)
	case http.MethodPost:
		switch cleaner.Line(request.Form.Get("method")) {
		case http.MethodPatch:
			_, err := editHistory(user.ID, media, owner, post, categories)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				log.Println(err)
				return
			}

			http.Redirect(writer, request, request.Referer(), http.StatusTemporaryRedirect)
		case http.MethodDelete:
			if file == "" {
				http.Error(writer, "file URL must be valid", http.StatusBadRequest)
				log.Println(err)
				return
			}

			history, err := deleteFileFromHistory(user.ID, owner, media, post, file)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				log.Println(err)
				return
			}

			redirectURL := request.Referer()
			if len(history.URLs) == 0 {
				URL, _ := url.Parse(redirectURL)
				query := URL.Query()
				query.Del("post")
				URL.RawQuery = query.Encode()
				redirectURL = URL.String()
			}
			http.Redirect(writer, request, redirectURL, http.StatusTemporaryRedirect)
		default:
			http.Error(writer, "request method is not recognized", http.StatusBadRequest)
			return
		}
	default:
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func filterHistories(U_ID primitive.ObjectID, media, owner, post string, categories []string) ([]db.History, error) {
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

func editHistory(U_ID primitive.ObjectID, media, owner, post string, categories []string) (db.History, error) {
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
	err := db.Histories.FindOneAndUpdate(context.Background(), filter, update, db.UpdateOptions).Decode(&history)
	return history, err
}

func deleteFileFromHistory(U_ID primitive.ObjectID, owner, media, post, file string) (db.History, error) {
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

	var history db.History
	if err := db.Histories.FindOneAndUpdate(context.Background(), filter, update, db.UpdateOptions).Decode(&history); err != nil {
		return db.History{}, err
	}

	if err := StorageHandler.Delete(media, owner, filepath.Base(file)); err != nil {
		return db.History{}, err
	}

	if len(history.URLs) == 0 {
		delete(filter, "urls")
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
