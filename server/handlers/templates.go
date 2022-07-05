package handlers

import (
	"html/template"
	"path/filepath"
	"strings"
)

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
		filepath.Join("templates", "authentication.html"),
		filepath.Join("templates", "history.html"),
		filepath.Join("templates", "histories.html"),
	))
)
