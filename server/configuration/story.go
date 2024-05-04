package configuration

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/AppleGamer22/raker/server/cleaner"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/AppleGamer22/raker/shared"
	"github.com/AppleGamer22/raker/shared/types"
	"github.com/charmbracelet/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (server *RakerServer) story(request *http.Request) (db.User, db.History, []error) {
	user := request.Context().Value(authenticatedUserKey).(db.User)

	if err := request.ParseForm(); err != nil {
		return db.User{}, db.History{}, []error{err}
	}

	history := db.History{
		Type: types.Story,
	}

	historyID := cleaner.Line(request.Form.Get("post"))
	owner := cleaner.Line(request.Form.Get("owner"))
	var errs []error

	if historyID != "" {
		filter := bson.M{
			"post":  historyID,
			"owner": owner,
			"type":  types.Story,
		}

		if err := server.Histories.FindOne(context.Background(), filter).Decode(&history); err == nil {
			return db.User{}, db.History{}, []error{err}
		}
	} else if owner == "" {
		return db.User{}, history, []error{}
	}

	instagram := shared.NewInstagram(user.Instagram.FBSR, user.Instagram.SessionID, user.Instagram.UserID)
	URLs, username, err := instagram.Reels(owner, false)
	if err != nil {
		return db.User{}, history, []error{err}
	}

	filter := bson.M{
		"type":  types.Story,
		"owner": username,
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
		filter["urls"] = fileName
		if count, err := server.Histories.CountDocuments(context.Background(), filter); err != nil || count > 0 {
			continue
		}

		localURLs = append(localURLs, fileName)
		newURLs = append(newURLs, urlString)
	}

	URLs = newURLs
	localURLs, saveErrors := StorageHandler.SaveBundle(user, types.Story, username, localURLs, URLs, []*http.Cookie{})
	errs = append(errs, saveErrors...)

	if len(localURLs) > 0 {
		historyID = primitive.NewObjectID().Hex()
		history = db.History{
			ID:    historyID,
			U_ID:  user.ID.Hex(),
			URLs:  localURLs,
			Type:  types.Story,
			Owner: username,
			Post:  historyID,
			Date:  time.Now(),
		}

		if _, err := server.Histories.InsertOne(context.Background(), history); err != nil {
			errs = append(errs, err)
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
