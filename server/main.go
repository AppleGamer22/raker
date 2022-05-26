package main

import (
	"log"
	"net/http"

	"github.com/AppleGamer22/rake/server/handlers"
)

func main() {
	http.HandleFunc("/instagram", handlers.Instagram)
	http.HandleFunc("/story", handlers.Story)
	http.HandleFunc("/tiktok", handlers.TikTok)
	http.HandleFunc("/vsco", handlers.VSCO)
	http.HandleFunc("/auth", handlers.Authentication)
	http.HandleFunc("/storage", handlers.Storage)
	log.Fatal(http.ListenAndServe(":4200", nil))
}
