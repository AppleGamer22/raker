package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/AppleGamer22/raker/server/authenticator"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/AppleGamer22/raker/server/handlers"
	"github.com/AppleGamer22/raker/shared"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func main() {
	if err1 := viper.ReadInConfig(); err1 != nil {
		if _, err := os.Stat("/.dockerenv"); err != nil {
			log.Error(err1)
		}
	}

	if err := viper.Unmarshal(&configuration); err != nil {
		log.Fatal(err)
	}

	if configuration.Secret == "" && !viper.IsSet("secret") {
		log.Fatal("A JWT secret must be set via a config file or an environment variable")
	}
	handlers.Authenticator = authenticator.New(configuration.Secret)
	client, err := db.Connect(configuration.URI, configuration.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	log.Infof("raker %s %s (%s/%s)", shared.Version, shared.Hash, runtime.GOOS, runtime.GOARCH)
	log.Infof("Storage path: %s", configuration.Storage)
	if configuration.Directories {
		log.Info("allowing directory listing")
	}
	log.Infof("MongoDB database URI: %s", configuration.URI)
	log.Infof("MongoDB database: %s", configuration.Database)
	log.Infof("Server is listening at http://localhost:%d", configuration.Port)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth/sign_up/instagram", handlers.InstagramSignUp)
	mux.HandleFunc("/api/auth/sign_in/instagram", handlers.InstagramSignIn)
	mux.HandleFunc("/api/auth/update/instagram", handlers.InstagramUpdateCredentials)
	mux.HandleFunc("/api/auth/sign_out/instagram", handlers.InstagramSignOut)
	mux.HandleFunc("/api/categories", handlers.Categories)
	mux.HandleFunc("/api/history", handlers.History)
	mux.HandleFunc("/api/info", handlers.Information)
	mux.Handle("/api/storage/", http.StripPrefix("/api/storage", handlers.NewStorageHandler(configuration.Storage, configuration.Directories)))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.Handle("/favicon.ico", http.RedirectHandler("/assets/icons/favicon.ico", http.StatusPermanentRedirect))
	mux.Handle("/robots.txt", http.RedirectHandler("/assets/robots.txt", http.StatusPermanentRedirect))

	mux.HandleFunc("/", handlers.AuthenticationPage)
	mux.HandleFunc("/history", handlers.HistoryPage)
	mux.HandleFunc("/instagram", handlers.InstagramPage)
	mux.HandleFunc("/highlight", handlers.HighlightPage)
	mux.HandleFunc("/story", handlers.StoryPage)
	mux.HandleFunc("/tiktok", handlers.TikTokPage)
	mux.HandleFunc("/vsco", handlers.VSCOPage)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", configuration.Port),
		Handler: handlers.NewLoggerMiddleware(mux),
		ErrorLog: log.Default().StandardLog(log.StandardLogOptions{
			ForceLevel: log.ErrorLevel,
		}),
	}

	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error(err)
			signals <- os.Interrupt
		}
	}()

	<-signals
	fmt.Print("\r")
	log.Warn("shutting down server...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Warn(err)
	}
}
