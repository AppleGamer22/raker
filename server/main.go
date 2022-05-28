package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AppleGamer22/rake/server/handlers"
	"github.com/spf13/viper"
)

func main() {
	viper.SetEnvPrefix("rake")
	viper.AutomaticEnv()
	viper.SetConfigName(".rake")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	if conf.Secret == "" && !viper.IsSet("secret") {
		log.Fatal("A JWT secret must be set via a config file or an environment variable")
	}

	log.Printf("Storage path: %s\n", conf.Storage)
	log.Printf("Users path: %s\n", conf.Users)
	log.Printf("MongoDB database URL: %s", conf.Database)
	log.Printf("Server is listening at http://localhost:%d\n", conf.Port)

	http.HandleFunc("/api/auth", handlers.Authentication)
	http.HandleFunc("/api/history", handlers.History)
	http.HandleFunc("/api/api/info", handlers.Information)
	http.HandleFunc("/api/instagram", handlers.Instagram)
	http.HandleFunc("/api/storage", handlers.Storage)
	http.HandleFunc("/api/story", handlers.Story)
	http.HandleFunc("/api/tiktok", handlers.TikTok)
	http.HandleFunc("/api/vsco", handlers.VSCO)

	http.HandleFunc("/auth", handlers.AuthenticationPage)
	http.HandleFunc("/history", handlers.HistoryPage)
	http.HandleFunc("/instagram", handlers.InstagramPage)
	http.HandleFunc("/story", handlers.StoryPage)
	http.HandleFunc("/tiktok", handlers.TikTokPage)

	fs := http.FileServer(http.Dir(conf.Storage))
	http.Handle("/api/storage/", http.StripPrefix("/api/storage/", fs))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil))
}
