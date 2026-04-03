package configuration

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/AppleGamer22/raker/server/cleaner"
	"github.com/AppleGamer22/raker/server/db"
	old "github.com/AppleGamer22/raker/server/db/mongo"
	"github.com/AppleGamer22/raker/shared"
	"github.com/AppleGamer22/raker/shared/types"
	"github.com/AppleGamer22/raker/templates"
	"github.com/charmbracelet/log"
)

func (server *RakerServer) History(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value(authenticatedUserKey).(db.User)

	if err := request.ParseForm(); err != nil {
		http.Error(writer, "failed to read request form", http.StatusBadRequest)
		return
	}
	media := cleaner.Line(request.Form.Get("media"))
	owner := cleaner.Line(request.Form.Get("owner"))
	post := cleaner.Line(request.Form.Get("post"))
	file := cleaner.Line(request.Form.Get("remove"))

	categories := make([]string, 0, len(user.Categories))
	for _, category := range user.Categories {
		if cleaner.Line(request.Form.Get(category)) == category && !types.ValidMediaType(category) {
			categories = append(categories, category)
		}
	}
	sort.Strings(categories)

	if media == "" || owner == "" || post == "" {
		http.Error(writer, "media type, owner & post must be valid", http.StatusBadRequest)
		return
	}

	switch request.Method {
	case http.MethodPatch:
		history, err := server.editHistory(user.Username, media, owner, post, categories)
		if err != nil {
			log.Error(err)
		}

		historyDisplay := old.HistoryDisplay{
			History:            history,
			Errors:             []error{err},
			SelectedCategories: user.SelectedCategories(history.Categories),
		}

		if err := templates.Templates.ExecuteTemplate(writer, "edit_categories.html", historyDisplay); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			log.Error(err)
		}
	case http.MethodDelete:
		if file == "" {
			http.Error(writer, "file URL must be valid", http.StatusBadRequest)
			return
		}

		history, err := server.deleteFileFromHistory(user, owner, media, post, file)
		if err != nil {
			log.Error(err)
		}

		// redirectURL := request.Referer()
		// URL, _ := url.Parse(redirectURL)
		// query := URL.Query()
		// if len(history.URLs) == 0 {
		// 	query.Del("post")
		// 	query.Del("owner")
		// } else if history.Type == types.Story && !query.Has("post") {
		// 	query.Set("post", history.Post)
		// }
		// URL.RawQuery = query.Encode()
		// redirectURL = URL.String()
		// http.Redirect(writer, request, redirectURL, http.StatusTemporaryRedirect)
		historyDisplay := old.HistoryDisplay{
			History:            history,
			Errors:             []error{err},
			SelectedCategories: user.SelectedCategories(history.Categories),
		}
		if err := templates.Templates.ExecuteTemplate(writer, "history_result.html", historyDisplay); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			log.Error(err)
		}
	default:
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func (server *RakerServer) filterHistories(user db.User, owner string, categories, mediaTypes []string, page int, exclusive bool) ([][]db.History, int, int, int, error) {
	for _, mediaType := range mediaTypes {
		if !types.ValidMediaType(mediaType) {
			return [][]db.History{}, 0, 0, 0, errors.New("media type must be valid")
		}
	}

	count, err := server.DBClient.HistoryCount(context.Background())
	if err != nil {
		return [][]db.History{}, 0, 0, 0, err
	}

	pages := int(math.Ceil(float64(count) / 30.0))
	if pages == 0 {
		page = 1
		pages = 1
	} else if page > pages {
		page = pages
	}

	if exclusive {
		sort.Strings(categories)
	}

	postTypes := make([]db.PostType, 0, len(mediaTypes))
	for _, mediaType := range mediaTypes {
		postTypes = append(postTypes, db.PostType(mediaType))
	}

	histories, err := server.DBClient.HistoryGetPage(context.Background(), db.HistoryGetPageParams{
		PostTypes:  postTypes,
		Exclusive:  exclusive,
		Categories: categories,
		PostOwner:  owner,
		Username:   user.Username,
		Page:       int32((page - 1) * 30),
		PageSize:   30,
	})
	if err != nil {
		return [][]db.History{}, 0, 0, 0, err
	}

	matrix := make([][]db.History, 0, int(math.Ceil(float64(len(histories))/3.0)))
	for i := 0; i < len(histories); i += 3 {
		end := i + 3
		if end > len(histories) {
			end = len(histories)
		}
		matrix = append(matrix, histories[i:end])
	}

	return matrix, page, pages, int(count), err
}

func (server *RakerServer) editHistory(username, media, owner, post string, categories []string) (db.History, error) {

	var history db.History
	err := server.DBClient.HistoryUpdateCategories(context.Background(), db.HistoryUpdateCategoriesParams{
		Categories: categories,
		Type:       db.PostType(media),
		Post:       post,
		Username:   username,
	})
	return history, err
}

func (server *RakerServer) deleteFileFromHistory(user db.User, owner, media, post, file string) (db.History, error) {
	history, err := server.DBClient.UpdateHistoryRemoveFile(context.Background(), db.UpdateHistoryRemoveFileParams{
		File:     file,
		Type:     db.PostType(media),
		Post:     post,
		Username: user.Username,
	})
	if err != nil {
		return db.History{}, err
	}

	if err := StorageHandler.Delete(user, media, owner, filepath.Base(file)); err != nil {
		return db.History{}, err
	}

	if len(history.Files) == 0 {
		err := server.DBClient.HistoryRemove(context.Background(), db.HistoryRemoveParams{
			Type:      db.PostType(media),
			Post:      post,
			Username:  user.Username,
			PostOwner: owner,
		})
		if err != nil {
			return db.History{}, err
		}
		return db.History{}, nil
	}

	return history, nil
}

func historyHTML(user db.User, history db.History, serverErrors []error, writer http.ResponseWriter) {
	historyDisplay := old.HistoryDisplay{
		History:            history,
		Errors:             serverErrors,
		Version:            shared.Version,
		SelectedCategories: user.SelectedCategories(history.Categories),
	}

	if err := templates.Templates.ExecuteTemplate(writer, "history.html", historyDisplay); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
}

func (server *RakerServer) LocationExif(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value(authenticatedUserKey).(db.User)
	if request.PathValue("user") != user.Username {
		http.Error(writer, "unathorised", http.StatusUnauthorized)
		return
	}
	latitude, longitude := StorageHandler.LocationEXIF(user, "vsco", request.PathValue("owner"), request.PathValue("file"))
	if latitude == 0 && longitude == 0 {
		writer.WriteHeader(http.StatusNotFound)
		writer.Header().Set("Content-Type", "text/html")
		fmt.Fprint(writer, `<html><script>window.close()</script></html>`)
		return
	}
	mapsURL := fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%f,%f", latitude, longitude)
	http.Redirect(writer, request, mapsURL, http.StatusTemporaryRedirect)
}

func (server *RakerServer) HistoryPage(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value(authenticatedUserKey).(db.User)

	if err := request.ParseForm(); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	owner := cleaner.Line(request.Form.Get("owner"))
	page, err := strconv.Atoi(cleaner.Line(request.Form.Get("page")))
	if err != nil {
		page = 1
	}

	categories := make([]string, 0, len(user.Categories))
	for _, category := range user.Categories {
		if cleaner.Line(request.Form.Get(category)) == category && !types.ValidMediaType(category) {
			categories = append(categories, category)
		}
	}

	mediaTypes := make([]string, 0, 5)
	for _, mediaType := range types.MediaTypes {
		if cleaner.Line(request.Form.Get(mediaType)) == mediaType && types.ValidMediaType(mediaType) {
			mediaTypes = append(mediaTypes, mediaType)
		}
	}
	if len(mediaTypes) == 0 {
		mediaTypes = types.MediaTypes
	}

	exclusive := cleaner.Line(request.Form.Get("exclusive")) == "exclusive"
	histories, page, pages, count, err := server.filterHistories(user, owner, categories, mediaTypes, page, exclusive)
	if err != nil {
		log.Error(err)
	}

	historiesDisplay := old.HistoriesDisplay{
		PostOwner:  owner,
		Categories: user.SelectedCategories(categories),
		Types:      db.SelectedMediaTypes(mediaTypes),
		Histories:  histories,
		Exclusive:  exclusive,
		Page:       page,
		Pages:      pages,
		Count:      count,
		Version:    shared.Version,
		Error:      err,
	}

	if err := templates.Templates.ExecuteTemplate(writer, "histories.html", historiesDisplay); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
}
