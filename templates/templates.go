package templates

import (
	"embed"
	"html/template"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	imagePathRegex = regexp.MustCompile("\\.(jpg)|(jpeg)|(webp)|(heic)$")
	videoPathRegex = regexp.MustCompile("\\.(mp4)|(webm)$")
	funcs          = template.FuncMap{
		"hasSuffix": strings.HasSuffix,
		"base":      filepath.Base,
		"add": func(a, b int) int {
			return a + b
		},
		"date": func(date time.Time) string {
			return date.Local().Format(time.RFC3339)
		},
		"isImagePath": func(path string) bool {
			return imagePathRegex.MatchString(path)
		},
		"isVideoPath": func(path string) bool {
			return videoPathRegex.MatchString(path)
		},
	}
	//go:embed *.html
	fs        embed.FS
	Templates = template.Must(template.New("").Funcs(funcs).ParseFS(fs, "*.html"))
)
