package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/AppleGamer22/raker/server/cleaner"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/AppleGamer22/raker/shared"
	"github.com/AppleGamer22/raker/shared/types"
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
		Type: types.Highlight,
	}

	highlightID := cleaner.Line(request.Form.Get("post"))
	errs := []error{}

	if highlightID != "" {
		filter := bson.M{
			"post": highlightID,
			"type": types.Highlight,
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
				localURLs = append(localURLs, fileName)
			}

			localURLs, saveErrors := StorageHandler.SaveBundle(user, types.Highlight, username, localURLs, URLs)
			errs = append(errs, saveErrors...)
			for _, err := range saveErrors {
				log.Println(err)
			}

			if len(localURLs) > 0 {
				history = db.History{
					ID:    primitive.NewObjectID().Hex(),
					U_ID:  user.ID.Hex(),
					URLs:  localURLs,
					Type:  types.Highlight,
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
