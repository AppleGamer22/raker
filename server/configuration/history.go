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
	db "github.com/AppleGamer22/raker/server/db/mongo"
	"github.com/AppleGamer22/raker/shared"
	"github.com/AppleGamer22/raker/shared/types"
	"github.com/AppleGamer22/raker/templates"
	"github.com/charmbracelet/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		history, err := server.editHistory(user.ID, media, owner, post, categories)
		if err != nil {
			log.Error(err)
		}

		historyDisplay := db.HistoryDisplay{
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
		historyDisplay := db.HistoryDisplay{
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

	filter := bson.M{"U_ID": user.ID.Hex()}

	if len(categories) != 0 {
		equal := len(categories) == len(user.Categories)
		if equal {
			for i := 0; i < len(categories); i++ {
				if categories[i] != user.Categories[i] {
					equal = false
					break
				}
			}
		}
		if exclusive {
			sort.Strings(categories)
			filter["categories"] = categories
		} else {
			if !equal {
				filter["categories"] = bson.M{
					"$in": categories,
				}
			}
		}
	} else {
		filter["$or"] = bson.A{
			bson.M{
				"categories": bson.M{
					"$size": 0,
				},
			},
			bson.M{
				"categories": bson.M{
					"$exists": false,
				},
			},
			bson.M{
				"categories": nil,
			},
		}
	}

	if len(mediaTypes) > 0 {
		filter["type"] = bson.M{
			"$in": mediaTypes,
		}
	}

	if owner != "" {
		filter["owner"] = bson.M{"$regex": primitive.Regex{
			Pattern: owner,
			Options: "i",
		}}
	}

	count, err := server.Histories.CountDocuments(context.Background(), filter)
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

	paginationOptions := options.Find()
	paginationOptions.SetSkip(int64((page - 1) * 30))
	paginationOptions.SetLimit(int64(30))
	paginationOptions.SetSort(bson.D{{Key: "date", Value: -1}})

	cursor, err := server.Histories.Find(context.Background(), filter, paginationOptions)
	if err != nil {
		return [][]db.History{}, 0, 0, 0, err
	}
	defer cursor.Close(context.Background())

	matrix := [][]db.History{}
	row := make([]db.History, 0, 3)
	for cursor.Next(context.Background()) {
		if len(row) == 3 {
			matrix = append(matrix, row)
			row = make([]db.History, 0, 3)
		}
		var history db.History
		if err = cursor.Decode(&history); err != nil {
			break
		}
		row = append(row, history)
	}
	if len(row) > 0 {
		matrix = append(matrix, row)
	}

	return matrix, page, pages, int(count), err
}

func (server *RakerServer) editHistory(U_ID primitive.ObjectID, media, owner, post string, categories []string) (db.History, error) {
	filter := bson.M{
		"U_ID":  U_ID.Hex(),
		"type":  media,
		"owner": owner,
		"post":  post,
	}

	update := bson.M{
		"$set": bson.M{
			"categories": categories,
		},
	}

	var history db.History
	err := server.Histories.FindOneAndUpdate(context.Background(), filter, update, db.UpdateOption).Decode(&history)
	return history, err
}

func (server *RakerServer) deleteFileFromHistory(user db.User, owner, media, post, file string) (db.History, error) {
	filter := bson.M{
		"U_ID": user.ID.Hex(),
		"urls": file,
		"post": post,
	}

	update := bson.M{
		"$pull": bson.M{
			"urls": file,
		},
	}

	var history db.History

	if err := server.Histories.FindOneAndUpdate(context.Background(), filter, update, db.UpdateOption).Decode(&history); err != nil {
		return db.History{}, err
	}

	if err := StorageHandler.Delete(user, media, owner, filepath.Base(file)); err != nil {
		return db.History{}, err
	}

	if len(history.URLs) == 0 {
		delete(filter, "urls")
		result, err := server.Histories.DeleteOne(context.Background(), filter)
		if err != nil {
			return db.History{}, err
		} else if result.DeletedCount == 0 {
			return db.History{}, errors.New("no histories were found")
		} else {
			return db.History{}, nil
		}
	}

	return history, nil
}

func historyHTML(user db.User, history db.History, serverErrors []error, writer http.ResponseWriter) {
	historyDisplay := db.HistoryDisplay{
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
	if request.PathValue("user") != user.ID.Hex() {
		http.Error(writer, "unathorised", http.StatusUnauthorized)
		return
	}
	latitude, longitude := StorageHandler.LocationEXIF(user, "vsco", request.PathValue("owner"), request.PathValue("file"))
	if latitude == 0 && longitude == 0 {
		writer.Header().Set("Content-Type", "text/html")
		fmt.Fprint(writer, `<html><script>window.close()</script></html>`)
		writer.WriteHeader(http.StatusNotFound)
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

	historiesDisplay := db.HistoriesDisplay{
		Owner:      owner,
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
