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

func (server *RakerServer) instagram(request *http.Request) (db.User, db.History, []error) {
	user := request.Context().Value(authenticatedUserKey).(db.User)

	if err := request.ParseForm(); err != nil {
		return db.User{}, db.History{}, []error{err}
	}

	history := db.History{
		Type: types.Instagram,
	}

	post := cleaner.Line(request.Form.Get("post"))
	incognito := cleaner.Line(request.Form.Get("incognito")) == "incognito"
	var errs []error

	if post != "" {
		filter := bson.M{
			"post": post,
			"type": types.Instagram,
		}
		if err := server.Histories.FindOne(context.Background(), filter).Decode(&history); err != nil {
			instagram := shared.NewInstagram(user.Instagram.FBSR, user.Instagram.SessionID, user.Instagram.UserID)
			var (
				username string
				URLs     []string
			)
			if incognito {
				URLs, username, _, err = shared.InstagramIncognito(post)
			} else {
				URLs, username, err = instagram.Post(post)
			}
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
				fileName := fmt.Sprintf("%s_%s", post, path.Base(URL.Path))
				localURLs = append(localURLs, fileName)
			}

			localURLs, saveErrors := StorageHandler.SaveBundle(user, types.Instagram, username, localURLs, URLs, []*http.Cookie{})
			errs = append(errs, saveErrors...)

			if len(localURLs) > 0 {
				history = db.History{
					ID:        primitive.NewObjectID().Hex(),
					U_ID:      user.ID.Hex(),
					URLs:      localURLs,
					Type:      types.Instagram,
					Owner:     username,
					Post:      post,
					Incognito: incognito,
					Date:      time.Now(),
				}

				if _, err := server.Histories.InsertOne(context.Background(), history); err != nil {
					errs = append(errs, err)
				}
			}
		}
	}
	return user, history, errs
}

func (server *RakerServer) InstagramPage(writer http.ResponseWriter, request *http.Request) {
	user, history, errs := server.instagram(request)
	if len(errs) > 0 {
		writer.WriteHeader(http.StatusBadRequest)
		for _, err := range errs {
			log.Error(err)
		}
	}
	historyHTML(user, history, errs, writer)
}

func (server *RakerServer) InstagramResult(writer http.ResponseWriter, request *http.Request) {
	user, history, errs := server.instagram(request)
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
