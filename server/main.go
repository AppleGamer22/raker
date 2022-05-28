package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AppleGamer22/rake/server/handlers"
	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault("databaseURL", "mongodb://localhost:27017/rake")
	viper.SetDefault("storagePath", ".")
	viper.SetDefault("usersPath", ".")
	viper.SetDefault("port", 4200)

	viper.SetEnvPrefix("rake")
	viper.SetConfigName("rake")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	if !viper.IsSet("secret") {
		log.Fatal("A JWT secret must be set via a config file or an environment variable")
	}

	log.Printf("MongoDB database URL: %s", conf.database)
	log.Printf("Server is listening at TCP port %d\n", conf.port)
	log.Printf("Storage path: %s\n", conf.storage)
	log.Printf("Users path: %s\n", conf.users)
	log.Println(conf.secret)
	os.Exit(0)
	http.HandleFunc("/instagram", handlers.Instagram)
	http.HandleFunc("/story", handlers.Story)
	http.HandleFunc("/tiktok", handlers.TikTok)
	http.HandleFunc("/vsco", handlers.VSCO)
	http.HandleFunc("/auth", handlers.Authentication)
	http.HandleFunc("/storage", handlers.Storage)
	http.HandleFunc("/version", handlers.Version)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", conf.port), nil))
}
