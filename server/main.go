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
	log.Printf("Server is listening at TCP port %d\n", conf.Port)

	http.HandleFunc("/instagram", handlers.Instagram)
	http.HandleFunc("/story", handlers.Story)
	http.HandleFunc("/tiktok", handlers.TikTok)
	http.HandleFunc("/vsco", handlers.VSCO)
	http.HandleFunc("/auth", handlers.Authentication)
	http.HandleFunc("/storage", handlers.Storage)
	http.HandleFunc("/version", handlers.Version)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil))
}
