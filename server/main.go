package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/AppleGamer22/rake/server/db"
	"github.com/AppleGamer22/rake/server/handlers"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(conf.URI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	database := client.Database(conf.Database, options.Database())
	db.Histories = *database.Collection("histories")
	db.Users = *database.Collection("users")

	log.Printf("Storage path: %s\n", conf.Storage)
	log.Printf("Users path: %s\n", conf.Users)
	log.Printf("MongoDB database URL: %s/%s", conf.URI, conf.Database)
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

	if err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil); err != nil {
		log.Fatal(err)
	}
}
