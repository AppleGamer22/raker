package configuration

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/AppleGamer22/raker/server/cleaner"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/AppleGamer22/raker/shared"
	"github.com/AppleGamer22/raker/shared/types"
	"github.com/charmbracelet/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (server *RakerServer) tiktok(request *http.Request) (db.User, db.History, []error) {
	user := request.Context().Value(authenticatedUserKey).(db.User)

	if err := request.ParseForm(); err != nil {
		return db.User{}, db.History{}, []error{err}
	}

	history := db.History{
		Type: types.TikTok,
	}
	owner := cleaner.Line(request.Form.Get("owner"))
	post := cleaner.Line(request.Form.Get("post"))
	incognito := cleaner.Line(request.Form.Get("incognito")) == "incognito"
	var errs []error

	if post != "" {
		filter := bson.M{
			"post": post,
			"type": types.TikTok,
		}
		if err := server.Histories.FindOne(context.Background(), filter).Decode(&history); err != nil {
			tiktok := shared.NewTikTok(user.TikTok.SessionID, user.TikTok.SessionIDGuard, user.TikTok.ChainToken)
			URLs, username, cookies, err := tiktok.Post(owner, post, incognito)
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
				if URL.Query().Get("mime_type") == "video_mp4" {
					localURLs = append(localURLs, fmt.Sprintf("%s.mp4", post))
					continue
				}
				fileName := fmt.Sprintf("%s_%s", post, path.Base(URL.Path))
				if !strings.HasSuffix(fileName, ".jpeg") {
					fileName = fmt.Sprintf("%s.jpeg", fileName)
				}
				localURLs = append(localURLs, fileName)
			}

			localURLs, saveErrors := StorageHandler.SaveBundle(user, types.TikTok, username, localURLs, URLs, cookies)
			errs = append(errs, saveErrors...)

			if len(localURLs) > 0 {
				history = db.History{
					ID:    primitive.NewObjectID().Hex(),
					U_ID:  user.ID.Hex(),
					URLs:  localURLs,
					Type:  types.TikTok,
					Owner: username,
					Post:  post,
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

func (server *RakerServer) TikTokPage(writer http.ResponseWriter, request *http.Request) {
	user, history, errs := server.tiktok(request)
	if len(errs) > 0 {
		writer.WriteHeader(http.StatusBadRequest)
		for _, err := range errs {
			log.Error(err)
		}
	}
	historyHTML(user, history, errs, writer)
}

func (server *RakerServer) TikTokResult(writer http.ResponseWriter, request *http.Request) {
	user, history, errs := server.tiktok(request)
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
	if err := templates.ExecuteTemplate(writer, "history_result.html", historyDisplay); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		log.Error(err)
	}
}
