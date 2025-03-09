package configuration

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/AppleGamer22/raker/server/cleaner"
	db "github.com/AppleGamer22/raker/server/db/mongo"
	"github.com/AppleGamer22/raker/shared"
	"github.com/AppleGamer22/raker/shared/types"
	"github.com/AppleGamer22/raker/templates"
	"github.com/charmbracelet/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (server *RakerServer) highlight(request *http.Request) (db.User, db.History, []error) {
	user := request.Context().Value(authenticatedUserKey).(db.User)

	if err := request.ParseForm(); err != nil {
		return db.User{}, db.History{}, []error{err}
	}

	history := db.History{
		Type: types.Highlight,
	}

	highlightID := cleaner.Line(request.Form.Get("post"))
	var errs []error

	if highlightID != "" {
		filter := bson.M{
			"post": highlightID,
			"type": types.Highlight,
		}

		if err := server.Histories.FindOne(context.Background(), filter).Decode(&history); err != nil {
			instagram := shared.NewInstagram(user.Instagram.FBSR, user.Instagram.SessionID, user.Instagram.UserID)
			URLs, username, err := instagram.Reels(highlightID, true)
			if err != nil {
				return db.User{}, db.History{}, []error{err}
			}

			localURLs := make([]string, 0, len(URLs))
			errs = make([]error, 0, len(URLs))
			for _, urlString := range URLs {
				URL, err := url.Parse(urlString)
				if err != nil {
					errs = append(errs, err)
					continue
				}
				fileName := fmt.Sprintf("%s_%s", highlightID, path.Base(URL.Path))
				localURLs = append(localURLs, fileName)
			}

			localURLs, saveErrors := StorageHandler.SaveBundle(user, types.Highlight, username, localURLs, URLs, []*http.Cookie{})
			errs = append(errs, saveErrors...)

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

				if _, err := server.Histories.InsertOne(context.Background(), history); err != nil {
					errs = append(errs, err)
				}
			}
		}
	}
	return user, history, errs
}

func (server *RakerServer) HighlightPage(writer http.ResponseWriter, request *http.Request) {
	user, history, errs := server.highlight(request)
	if len(errs) > 0 {
		writer.WriteHeader(http.StatusBadRequest)
		for _, err := range errs {
			log.Error(err)
		}
	}
	historyHTML(user, history, errs, writer)
}

func (server *RakerServer) HighlightResult(writer http.ResponseWriter, request *http.Request) {
	user, history, errs := server.highlight(request)
	if len(errs) > 0 {
		// writer.WriteHeader(http.StatusBadRequest)
		for _, err := range errs {
			log.Error(err)
		}
	}
	historyDisplay := db.HistoryDisplay{
		History:            history,
		Errors:             errs,
		SelectedCategories: user.SelectedCategories(history.Categories),
	}
	if err := templates.Templates.ExecuteTemplate(writer, "history_result.html", historyDisplay); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		log.Error(err)
	}
}
