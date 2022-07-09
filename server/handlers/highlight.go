package handlers

import (
	"context"
	"fmt"
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

func HighlightPage(writer http.ResponseWriter, request *http.Request) {
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
		Type: db.Highlight,
	}

	highlightID := cleaner.Line(request.Form.Get("post"))
	errs := []error{}

	if highlightID != "" {
		filter := bson.M{
			"post": highlightID,
			"type": db.Highlight,
		}

		if err := db.Histories.FindOne(context.Background(), filter).Decode(&history); err != nil {
			instagram := shared.NewInstagram(user.Instagram.FBSR, user.Instagram.SessionID, user.Instagram.UserID)
			URLs, username, err := instagram.Reels(highlightID, true)
			if err != nil {
				log.Println(err)
				historyHTML(user, history, []error{err}, writer)
				return
			}

			localURLs := make([]string, 0, len(URLs))
			for _, urlString := range URLs {
				URL, err := url.Parse(urlString)
				if err != nil {
					log.Println(err)
					errs = append(errs, err)
					continue
				}
				fileName := fmt.Sprintf("%s_%s", highlightID, path.Base(URL.Path))

				if err := StorageHandler.Save(user, db.Highlight, username, fileName, urlString); err != nil {
					log.Println(err)
					errs = append(errs, err)
					continue
				}
				localURLs = append(localURLs, fileName)
			}

			if len(localURLs) > 0 {
				history = db.History{
					ID:    primitive.NewObjectID().Hex(),
					U_ID:  user.ID.Hex(),
					URLs:  localURLs,
					Type:  db.Highlight,
					Owner: username,
					Post:  highlightID,
					Date:  time.Now(),
				}

				if _, err := db.Histories.InsertOne(context.Background(), history); err != nil {
					log.Println(err)
					errs = append(errs, err)
				}
			}
		}
	}

	historyHTML(user, history, errs, writer)
}
