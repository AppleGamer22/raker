package configuration

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/AppleGamer22/raker/server/cleaner"
	"github.com/AppleGamer22/raker/server/db"
	old "github.com/AppleGamer22/raker/server/db/mongo"
	"github.com/AppleGamer22/raker/shared"
	"github.com/AppleGamer22/raker/shared/types"
	"github.com/AppleGamer22/raker/templates"
	"github.com/charmbracelet/log"
)

func (server *RakerServer) tiktok(request *http.Request) (db.User, db.History, []error) {
	user := request.Context().Value(authenticatedUserKey).(db.User)

	if err := request.ParseForm(); err != nil {
		return db.User{}, db.History{}, []error{err}
	}

	history := db.History{
		PostType: types.TikTok,
	}
	owner := cleaner.Line(request.Form.Get("owner"))
	post := cleaner.Line(request.Form.Get("post"))
	incognito := cleaner.Line(request.Form.Get("incognito")) == "incognito"
	var errs []error

	if post != "" {
		retrievedHistory, err := server.DBClient.HistoryGet(context.Background(), db.HistoryGetParams{
			PostType: db.PostTypeTiktok,
			Post:     post,
			Username: user.Username,
		})
		if err == nil {
			history = retrievedHistory
		} else {
			tiktok := shared.NewTikTok(user.TiktokSessionID, user.TiktokSessionIDGuard)
			videoURLs, coverURLs, username, cookies, err := tiktok.Post(owner, post, incognito)
			if err != nil {
				return db.User{}, db.History{}, []error{err}
			}

			errs = make([]error, 0, len(coverURLs)+1)
			var videoURL string
			for i, urlString := range videoURLs {
				fileName := fmt.Sprintf("%s.mp4", post)
				if err := StorageHandler.Save(user, types.TikTok, username, fileName, urlString, cookies); err != nil {
					if i == len(videoURLs)-1 {
						errs = append(errs, err)
					}
					continue
				}
				videoURL = fileName
				break
			}

			localURLs := make([]string, 0, len(coverURLs)+1)
			for _, urlString := range coverURLs {
				URL, err := url.Parse(urlString)
				if err != nil {
					errs = append(errs, err)
					continue
				}
				fileName := fmt.Sprintf("%s_%s.jpeg", post, path.Base(URL.Path))
				localURLs = append(localURLs, fileName)
			}

			localURLs, saveErrors := StorageHandler.SaveBundle(user, types.TikTok, username, localURLs, coverURLs, cookies)
			errs = append(errs, saveErrors...)

			if len(videoURL) > 0 {
				localURLs = append([]string{videoURL}, localURLs...)
			}

			if len(localURLs) > 0 {
				addedHistory, err := server.DBClient.HistoryAdd(context.Background(), db.HistoryAddParams{
					Username:  user.Username,
					PostType:  db.PostTypeTiktok,
					PostOwner: username,
					Post:      post,
					Files:     localURLs,
				})
				if err != nil {
					return db.User{}, db.History{}, []error{err}
				}
				history = addedHistory
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
	historyDisplay := old.HistoryDisplay{
		History:            history,
		Errors:             errs,
		SelectedCategories: user.SelectedCategories(history.Categories),
	}
	if err := templates.Templates.ExecuteTemplate(writer, "history_result.html", historyDisplay); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		log.Error(err)
	}
}
