package handlers

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/AppleGamer22/rake/server/cleaner"
	"github.com/AppleGamer22/rake/server/db"
	"github.com/AppleGamer22/rake/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StoryPage(writer http.ResponseWriter, request *http.Request) {
	user, err := Verify(request)
	if err != nil {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	if err := request.ParseForm(); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	history := db.History{
		Type: db.Story,
	}

	historyID := cleaner.Line(request.Form.Get("post"))
	owner := cleaner.Line(request.Form.Get("owner"))
	errs := []error{}

	if historyID != "" {
		filter := bson.M{
			"post":  historyID,
			"owner": owner,
			"type":  db.Story,
		}

		if err := db.Histories.FindOne(context.Background(), filter).Decode(&history); err == nil {
			historyHTML(user, history, []error{}, writer)
			return
		}
	} else if owner == "" {
		historyHTML(user, history, errs, writer)
		return
	}

	instagram := shared.NewInstagram(user.Instagram.FBSR, user.Instagram.SessionID, user.Instagram.UserID)
	URLs, username, err := instagram.Reels(owner, false)
	if err != nil {
		log.Println(err)
		historyHTML(user, history, []error{err}, writer)
		return
	}

	filter := bson.M{
		"type":  db.Story,
		"owner": username,
	}
	localURLs := make([]string, 0, len(URLs))
	for _, urlString := range URLs {
		URL, err := url.Parse(urlString)
		if err != nil {
			log.Println(err)
			errs = append(errs, err)
			continue
		}

		fileName := path.Base(URL.Path)
		filter["urls"] = fileName
		if count, err := db.Histories.CountDocuments(context.Background(), filter); err != nil || count > 0 {
			continue
		}

		if err := StorageHandler.Save(user, db.Story, username, fileName, urlString); err != nil {
			log.Println(err)
			errs = append(errs, err)
			continue
		}

		localURLs = append(localURLs, fileName)
	}

	if len(localURLs) > 0 {
		historyID = primitive.NewObjectID().Hex()
		history = db.History{
			ID:    historyID,
			U_ID:  user.ID.Hex(),
			URLs:  localURLs,
			Type:  db.Story,
			Owner: username,
			Post:  historyID,
			Date:  time.Now(),
		}

		if _, err := db.Histories.InsertOne(context.Background(), history); err != nil {
			log.Println(err)
			errs = append(errs, err)
		}
	}

	historyHTML(user, history, errs, writer)
}
