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

func (server *RakerServer) instagram(request *http.Request) (db.User, db.History, []error) {
	user := request.Context().Value(authenticatedUserKey).(db.User)

	if err := request.ParseForm(); err != nil {
		return db.User{}, db.History{}, []error{err}
	}

	post := cleaner.Line(request.Form.Get("post"))
	incognito := cleaner.Line(request.Form.Get("incognito")) == "incognito"
	history := db.History{
		PostType:  db.PostTypeInstagram,
		Incognito: incognito,
	}
	var errs []error

	if post != "" {
		retrievedHistory, err := server.DBClient.HistoryGet(context.Background(), db.HistoryGetParams{
			PostType: db.PostTypeInstagram,
			Post:     post,
			Username: user.Username,
		})
		if err == nil {
			retrievedHistory.Incognito = incognito
			history = retrievedHistory
		} else {
			instagram := shared.NewInstagram(user.InstagramSessionID, user.InstagramUserID)
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
				addedHistory, err := server.DBClient.HistoryAdd(context.Background(), db.HistoryAddParams{
					Username:  user.Username,
					PostType:  db.PostTypeInstagram,
					PostOwner: username,
					Post:      post,
					Files:     localURLs,
				})
				if err != nil {
					errs = append(errs, err)
				} else {
					addedHistory.Incognito = incognito
					history = addedHistory
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
