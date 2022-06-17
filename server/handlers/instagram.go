package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func Instagram(writer http.ResponseWriter, request *http.Request) {

}

func InstagramPage(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("templates", "instagram.html"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(writer, "instagram", nil)
}
