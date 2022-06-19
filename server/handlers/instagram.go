package handlers

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/AppleGamer22/rake/server/db"
	"github.com/AppleGamer22/rake/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InstagramPage(writer http.ResponseWriter, request *http.Request) {
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

	var history db.History
	post := request.Form.Get("post")
	if post != "" {
		filter := bson.M{
			"post": post,
			"type": db.Instagram,
		}
		if err := db.Histories.FindOne(context.Background(), filter).Decode(&history); err != nil {
			instagram := shared.NewInstagram(user.FBSR, user.SessionID, user.AppID)
			URLs, username, err := instagram.Post(post)

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			localURLs := make([]string, 0, len(URLs))
			for _, urlString := range URLs {
				URL, err := url.Parse(urlString)
				if err != nil {
					log.Println(err)
					continue
				}
				fileName := fmt.Sprintf("%s_%s", post, path.Base(URL.Path))

				if err := StorageHandler.Save(db.Instagram, username, fileName, urlString); err != nil {
					log.Println(err)
					continue
				}
				localURLs = append(localURLs, fmt.Sprintf("storage/instagram/%s/%s", username, fileName))
			}

			history = db.History{
				ID:    primitive.NewObjectID().Hex(),
				U_ID:  user.ID.Hex(),
				URLs:  localURLs,
				Type:  db.Instagram,
				Owner: username,
				Post:  post,
				Date:  time.Now(),
			}

			if _, err := db.Histories.InsertOne(context.Background(), history); err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	funcs := template.FuncMap{
		"hasSuffix": strings.HasSuffix,
		"join":      strings.Join,
		"base":      filepath.Base,
	}
	tmpl, err := template.New("instagram.html").Funcs(funcs).ParseFiles(filepath.Join("templates", "instagram.html"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	historyDisplay := db.HistoryDisplay{
		History: history,
		AvailableCategories: func() map[string]bool {
			result := make(map[string]bool)
			for _, category := range user.Categories {
				result[category] = true
			}
			for _, category := range history.Categories {
				if _, ok := result[category]; ok {
					result[category] = false
				}
			}
			for category, c := range result {
				result[category] = !c
			}
			return result
		}(),
	}

	writer.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(writer, historyDisplay); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
