package handlers

import (
	"context"
	"errors"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/AppleGamer22/rake/server/cleaner"
	"github.com/AppleGamer22/rake/server/db"
	"github.com/AppleGamer22/rake/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func History(writer http.ResponseWriter, request *http.Request) {
	user, err := Verify(request)
	if err != nil {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		log.Println(err)
		return
	}

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
		if cleaner.Line(request.Form.Get(category)) == category && !db.ValidMediaType(category) {
			categories = append(categories, category)
		}
	}
	sort.Strings(categories)

	if media == "" || owner == "" || post == "" {
		http.Error(writer, "media type, owner & post must be valid", http.StatusBadRequest)
		return
	}

	switch request.Method {
	case http.MethodPost:
		switch cleaner.Line(request.Form.Get("method")) {
		case http.MethodPatch:
			_, err := editHistory(user.ID, media, owner, post, categories)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				log.Println(err)
				return
			}

			http.Redirect(writer, request, request.Referer(), http.StatusTemporaryRedirect)
		case http.MethodDelete:
			if file == "" {
				http.Error(writer, "file URL must be valid", http.StatusBadRequest)
				log.Println(err)
				return
			}

			history, err := deleteFileFromHistory(user, owner, media, post, file)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				log.Println(err)
				return
			}

			redirectURL := request.Referer()
			if len(history.URLs) == 0 {
				URL, _ := url.Parse(redirectURL)
				query := URL.Query()
				query.Del("post")
				URL.RawQuery = query.Encode()
				redirectURL = URL.String()
			}
			http.Redirect(writer, request, redirectURL, http.StatusTemporaryRedirect)
		default:
			http.Error(writer, "request method is not recognized", http.StatusBadRequest)
			return
		}
	default:
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func filterHistories(U_ID primitive.ObjectID, owner string, categories, mediaTypes []string, page int) ([][]db.History, int, error) {
	for _, mediaType := range mediaTypes {
		if !db.ValidMediaType(mediaType) {
			return [][]db.History{}, 0, errors.New("media type must be valid")
		}
	}

	filter := bson.M{"U_ID": U_ID.Hex()}

	if len(categories) != 0 {
		filter["categories"] = bson.M{
			"$in": categories,
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

	paginationOptions := options.Find().SetSkip(int64((page - 1) * 30)).SetLimit(int64(30)).SetSort(bson.M{"post": -1})
	cursor, err := db.Histories.Find(context.Background(), filter, paginationOptions)
	if err != nil {
		return [][]db.History{}, 0, err
	}
	var histories []db.History
	err = cursor.All(context.Background(), &histories)

	matrix := [][]db.History{}
	for i := len(histories) - 1; i > -1; i -= 3 {
		row := []db.History{}
		for j := 0; j < 3 && i-j > -1; j++ {
			row = append(row, histories[i-j])
			histories = histories[:i-j]
		}
		matrix = append(matrix, row)
	}

	count, _ := db.Histories.CountDocuments(context.Background(), filter)
	return matrix, int(count) / 30, err
}

func editHistory(U_ID primitive.ObjectID, media, owner, post string, categories []string) (db.History, error) {
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
	err := db.Histories.FindOneAndUpdate(context.Background(), filter, update, db.UpdateOption).Decode(&history)
	return history, err
}

func deleteFileFromHistory(user db.User, owner, media, post, file string) (db.History, error) {
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
	if err := db.Histories.FindOneAndUpdate(context.Background(), filter, update, db.UpdateOption).Decode(&history); err != nil {
		return db.History{}, err
	}

	if err := StorageHandler.Delete(user, media, owner, filepath.Base(file)); err != nil {
		return db.History{}, err
	}

	if len(history.URLs) == 0 {
		delete(filter, "urls")
		result, err := db.Histories.DeleteOne(context.Background(), filter)
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

var (
	funcs = template.FuncMap{
		"hasSuffix": strings.HasSuffix,
		// "join":      strings.Join,
		"base": filepath.Base,
		"add": func(a, b int) int {
			return a + b
		},
	}
	templates = template.Must(template.New("").Funcs(funcs).ParseFiles(
		filepath.Join("templates", "history.html"),
		filepath.Join("templates", "histories.html"),
	))
)

func historyHTML(user db.User, history db.History, serverErrors []error, writer http.ResponseWriter) {
	historyDisplay := db.HistoryDisplay{
		History:            history,
		Errors:             serverErrors,
		Version:            shared.Version,
		SelectedCategories: user.SelectedCategories(history.Categories),
	}

	writer.Header().Set("Content-Type", "text/html")
	if err := templates.ExecuteTemplate(writer, "history.html", historyDisplay); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func HistoryPage(writer http.ResponseWriter, request *http.Request) {
	user, err := Verify(request)
	if err != nil {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		log.Println(err)
		return
	}

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
		if cleaner.Line(request.Form.Get(category)) == category && !db.ValidMediaType(category) {
			categories = append(categories, category)
		}
	}

	mediaTypes := make([]string, 0, 5)
	for _, mediaType := range db.MediaTypes {
		if cleaner.Line(request.Form.Get(mediaType)) == mediaType && db.ValidMediaType(mediaType) {
			mediaTypes = append(mediaTypes, mediaType)
		}
	}
	if len(mediaTypes) == 0 {
		mediaTypes = db.MediaTypes
	}

	histories, pages, err := filterHistories(user.ID, owner, categories, mediaTypes, page)
	if err != nil {
		log.Println(err)
	}

	historiesDisplay := db.HistoriesDisplay{
		Owner:      owner,
		Categories: user.SelectedCategories(categories),
		Types:      db.SelectedMediaTypes(mediaTypes),
		Histories:  histories,
		Errors:     []error{err},
		Page:       page,
		Pages:      pages,
		Version:    shared.Version,
	}

	writer.Header().Set("Content-Type", "text/html")
	if err := templates.ExecuteTemplate(writer, "histories.html", historiesDisplay); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
