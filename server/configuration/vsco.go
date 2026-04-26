package configuration

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/AppleGamer22/raker/server/cleaner"
	"github.com/AppleGamer22/raker/server/db"
	old "github.com/AppleGamer22/raker/server/db/mongo"
	"github.com/AppleGamer22/raker/shared"
	"github.com/AppleGamer22/raker/shared/types"
	"github.com/AppleGamer22/raker/templates"
	"github.com/charmbracelet/log"
)

func (server *RakerServer) vsco(request *http.Request) (db.User, db.History, []error) {
	user := request.Context().Value(authenticatedUserKey).(db.User)

	if err := request.ParseForm(); err != nil {
		return db.User{}, db.History{}, []error{err}
	}

	history := db.History{
		PostType: types.VSCO,
	}
	owner := cleaner.Line(request.Form.Get("owner"))
	post := cleaner.Line(request.Form.Get("post"))
	var errs []error

	if post != "" {
		retrievedHistory, err := server.DBClient.HistoryGet(context.Background(), db.HistoryGetParams{
			PostType: db.PostTypeVsco,
			Post:     post,
			Username: user.Username,
		})
		if err == nil {
			history = retrievedHistory
		} else {
			URLs, username, cookies, err := shared.VSCO(owner, post)
			if err != nil {
				return db.User{}, db.History{}, []error{err}
			}

			localURLs := make([]string, 0, len(URLs))
			errs = make([]error, 0, len(URLs))
			for _, urlString := range URLs {

				URL, err := url.Parse(urlString)
				if err != nil {
					errs = append(errs, err)
				}
				if strings.Contains(urlString, ".ts") {
					localURLs = append(localURLs, fmt.Sprintf("%s.mp4", post))
				} else if strings.Contains(URL.Path, "/poster/private") {
					localURLs = append(localURLs, fmt.Sprintf("%s.jpg", post))
				} else {
					localURL := fmt.Sprintf("%s_%s", post, path.Base(URL.Path))
					if !strings.HasSuffix(localURL, ".jpg") && !strings.HasSuffix(localURL, ".jpeg") {
						localURL = fmt.Sprintf("%s.jpeg", localURL)
					}
					localURLs = append(localURLs, localURL)
				}
			}
			localURLs, saveErrors := StorageHandler.SaveBundle(user, types.VSCO, username, localURLs, URLs, cookies)
			errs = append(errs, saveErrors...)

			if len(localURLs) > 0 {
				addedHistory, err := server.DBClient.HistoryAdd(context.Background(), db.HistoryAddParams{
					Username:   user.Username,
					PostType:   db.PostTypeVsco,
					PostOwner:  username,
					Post:       post,
					Files:      localURLs,
					Categories: []string{},
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

func (server *RakerServer) VSCOPage(writer http.ResponseWriter, request *http.Request) {
	user, history, errs := server.vsco(request)
	if len(errs) > 0 {
		writer.WriteHeader(http.StatusBadRequest)
		for _, err := range errs {
			log.Error(err)
		}
	}
	historyHTML(user, history, nil, writer)
}

func (server *RakerServer) VSCOResult(writer http.ResponseWriter, request *http.Request) {
	user, history, errs := server.vsco(request)
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
