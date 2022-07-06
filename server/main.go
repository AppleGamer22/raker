package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/AppleGamer22/rake/server/authenticator"
	"github.com/AppleGamer22/rake/server/db"
	"github.com/AppleGamer22/rake/server/handlers"
	"github.com/spf13/viper"
)

func init() {
	viper.SetEnvPrefix("rake")
	viper.AutomaticEnv()
	viper.SetConfigName(".rake")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
}

func main() {
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

	client, err := db.Connect(conf.URI, conf.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	log.Printf("Storage path: %s\n", conf.Storage)
	if conf.Directories {
		log.Println("allowing directory listing")
	}
	log.Printf("MongoDB database URI: %s", conf.URI)
	log.Printf("MongoDB database: %s", conf.Database)
	log.Printf("Server is listening at http://localhost:%d\n", conf.Port)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth/sign_up/instagram", handlers.InstagramSignUp)
	mux.HandleFunc("/api/auth/sign_in/instagram", handlers.InstagramSignIn)
	mux.HandleFunc("/api/auth/update/instagram", handlers.InstagramUpdateCredentials)
	mux.HandleFunc("/api/auth/sign_out/instagram", handlers.InstagramSignOut)
	mux.HandleFunc("/api/history", handlers.History)
	mux.HandleFunc("/api/info", handlers.Information)
	mux.HandleFunc("/api/story", handlers.Story)
	mux.Handle("/api/storage/", http.StripPrefix("/api/storage", handlers.NewStorageHandler(conf.Storage, conf.Directories)))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	mux.HandleFunc("/", handlers.AuthenticationPage)
	mux.HandleFunc("/history", handlers.HistoryPage)
	mux.HandleFunc("/instagram", handlers.InstagramPage)
	mux.HandleFunc("/story", handlers.StoryPage)
	mux.HandleFunc("/tiktok", handlers.TikTokPage)
	mux.HandleFunc("/vsco", handlers.VSCOPage)

	server := http.Server{
		Addr:    fmt.Sprintf("localhost:%d", conf.Port),
		Handler: handlers.Log(mux),
	}

	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println(err)
			signals <- os.Interrupt
		}
	}()

	<-signals
	fmt.Print("\r")
	log.Println("shutting down server...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}
}
