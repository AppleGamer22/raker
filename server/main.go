package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/AppleGamer22/rake/server/authenticator"
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
	handlers.Authenticator = authenticator.New(conf.Secret)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(conf.URI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	database := client.Database(conf.Database, options.Database())
	db.Histories = *database.Collection("histories")
	db.Users = *database.Collection("users")

	log.Printf("Storage path: %s\n", conf.Storage)
	if conf.Directories {
		log.Println("allowing directory listing")
	}
	log.Printf("Users path: %s\n", conf.Users)
	log.Printf("MongoDB database URL: %s", path.Join(conf.URI, conf.Database))
	log.Printf("Server is listening at http://localhost:%d\n", conf.Port)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth", handlers.Authentication)
	mux.HandleFunc("/api/history", handlers.History)
	mux.HandleFunc("/api/api/info", handlers.Information)
	mux.HandleFunc("/api/instagram", handlers.Instagram)
	mux.HandleFunc("/api/story", handlers.Story)
	mux.HandleFunc("/api/tiktok", handlers.TikTok)
	mux.HandleFunc("/api/vsco", handlers.VSCO)
	mux.Handle("/api/storage/", http.StripPrefix("/api/storage", handlers.NewStorageHandler(conf.Storage, conf.Directories)))

	mux.HandleFunc("/auth", handlers.AuthenticationPage)
	mux.HandleFunc("/history", handlers.HistoryPage)
	mux.HandleFunc("/instagram", handlers.InstagramPage)
	mux.HandleFunc("/story", handlers.StoryPage)
	mux.HandleFunc("/tiktok", handlers.TikTokPage)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), handlers.Log(mux)); err != nil {
		log.Fatal(err)
	}
}
