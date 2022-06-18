package handlers

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/AppleGamer22/rake/server/db"
	"go.mongodb.org/mongo-driver/bson"
)

func Instagram(writer http.ResponseWriter, request *http.Request) {

}

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
			// http.Error(writer, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}

	funcs := template.FuncMap{
		"hasSuffix": strings.HasSuffix,
		"join":      strings.Join,
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
