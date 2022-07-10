package handlers

import (
	"html/template"
	"path/filepath"
	"strings"
	"time"
)

var (
	funcs = template.FuncMap{
		"hasSuffix": strings.HasSuffix,
		"base":      filepath.Base,
		"add": func(a, b int) int {
			return a + b
		},
		"date": func(date time.Time) string {
			return date.Local().Format(time.RFC3339)
		},
	}
	templates = template.Must(template.New("").Funcs(funcs).ParseFiles(
		filepath.Join("templates", "authentication.html"),
		filepath.Join("templates", "history.html"),
		filepath.Join("templates", "histories.html"),
	))
)
