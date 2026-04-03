package configuration

import (
	"context"
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
	"github.com/google/uuid"
)

func (server *RakerServer) story(request *http.Request) (db.User, db.History, []error) {
	user := request.Context().Value(authenticatedUserKey).(db.User)

	if err := request.ParseForm(); err != nil {
		return db.User{}, db.History{}, []error{err}
	}

	history := db.History{
		PostType: types.Story,
	}

	historyID := cleaner.Line(request.Form.Get("post"))
	owner := cleaner.Line(request.Form.Get("owner"))
	var errs []error

	if historyID != "" {
		retrievedHistory, err := server.DBClient.HistoryGetByOwner(context.Background(), db.HistoryGetByOwnerParams{
			PostType:  db.PostTypeStory,
			PostOwner: owner,
			Post:      historyID,
			Username:  user.Username,
		})
		if err == nil {
			// history log already exists
			return user, retrievedHistory, []error{}
		}
	} else if owner == "" {
		// empty input results in empty response
		return user, history, []error{}
	}

	instagram := shared.NewInstagram("", user.InstagramSessionID, user.InstagramUserID)
	URLs, username, err := instagram.Reels(owner, false)
	if err != nil {
		return user, history, []error{err}
	}

	newURLs := make([]string, 0, len(URLs))
	localURLs := make([]string, 0, len(URLs))
	errs = make([]error, 0, len(URLs))
	for _, urlString := range URLs {
		URL, err := url.Parse(urlString)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		fileName := path.Base(URL.Path)
		count, err := server.DBClient.HistoryCountByFile(context.Background(), db.HistoryCountByFileParams{
			PostType:  db.PostTypeStory,
			PostOwner: username,
			File:      fileName,
			Username:  user.Username,
		})
		if err != nil || count > 0 {
			continue
		}

		localURLs = append(localURLs, fileName)
		newURLs = append(newURLs, urlString)
	}

	URLs = newURLs
	localURLs, saveErrors := StorageHandler.SaveBundle(user, types.Story, username, localURLs, URLs, []*http.Cookie{})
	errs = append(errs, saveErrors...)

	if len(localURLs) > 0 {
		if historyID == "" {
			historyID = uuid.NewString()
		}
		addedHistory, err := server.DBClient.HistoryAdd(context.Background(), db.HistoryAddParams{
			Username:  user.Username,
			PostType:  db.PostTypeStory,
			PostOwner: username,
			Post:      historyID,
			Files:     localURLs,
		})
		if err != nil {
			errs = append(errs, err)
		} else {
			history = addedHistory
		}
	}

	return user, history, errs
}

func (server *RakerServer) StoryPage(writer http.ResponseWriter, request *http.Request) {
	user, history, errs := server.story(request)
	if len(errs) > 0 {
		writer.WriteHeader(http.StatusBadRequest)
		for _, err := range errs {
			log.Error(err)
		}
	}

	historyHTML(user, history, errs, writer)
}

func (server *RakerServer) StoryResult(writer http.ResponseWriter, request *http.Request) {
	user, history, errs := server.story(request)
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
