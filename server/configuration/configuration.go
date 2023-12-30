package configuration

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AppleGamer22/raker/server/authenticator"
	"github.com/AppleGamer22/raker/server/db"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type contextKey int

const authenticatedUserKey contextKey = 0

type Configuration struct {
	Secret      string
	URI         string
	Database    string
	Storage     string
	Directories bool
	Port        uint
}

type RakerServer struct {
	Configuration
	DBClient      *mongo.Client
	Users         *mongo.Collection
	Histories     *mongo.Collection
	Authenticator authenticator.Authenticator
	HTTPServer    http.Server
}

func NewRakerServer() (*RakerServer, error) {
	rakerServer := RakerServer{
		Configuration: Configuration{
			URI:         "mongodb://localhost:27017",
			Database:    "raker",
			Storage:     ".",
			Directories: false,
			Port:        4100,
		},
	}

	if err1 := viper.ReadInConfig(); err1 != nil {
		if _, err := os.Stat("/.dockerenv"); err != nil {
			log.Error(err1)
		}
	}

	if err := viper.Unmarshal(&rakerServer.Configuration); err != nil {
		log.Fatal(err)
	}

	if configuration.Secret == "" && !viper.IsSet("secret") {
		log.Fatal("A JWT secret must be set via a config file or an environment variable")
	}

	rakerServer.Authenticator = authenticator.New(configuration.Secret)

	dbClient, database, err := db.Connect(configuration.URI, configuration.Database)
	if err != nil {
		log.Fatal(err)
	}
	rakerServer.DBClient = dbClient
	// remember to defer client.Disconnet
	rakerServer.Histories = database.Collection("histories")
	rakerServer.Users = database.Collection("users")

	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth/sign_up/instagram", rakerServer.InstagramSignUp)
	mux.HandleFunc("/api/auth/sign_in/instagram", rakerServer.InstagramSignIn)
	mux.Handle("/api/auth/update/instagram", rakerServer.Verify(http.HandlerFunc(rakerServer.InstagramUpdateCredentials)))
	mux.Handle("/api/auth/sign_out/instagram", rakerServer.Verify(http.HandlerFunc(rakerServer.InstagramSignOut)))
	mux.Handle("/api/categories", rakerServer.Verify(http.HandlerFunc(rakerServer.Categories)))
	mux.Handle("/api/history", rakerServer.Verify(http.HandlerFunc(rakerServer.History)))
	// mux.HandleFunc("/api/info", rakerServer.Information)
	mux.Handle("/api/storage/", http.StripPrefix("/api/storage", rakerServer.Verify(NewStorageHandler(configuration.Storage, configuration.Directories))))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.Handle("/favicon.ico", http.RedirectHandler("/assets/icons/favicon.ico", http.StatusPermanentRedirect))
	mux.Handle("/robots.txt", http.RedirectHandler("/assets/robots.txt", http.StatusPermanentRedirect))

	mux.HandleFunc("/", rakerServer.AuthenticationPage)
	mux.Handle("/history", rakerServer.Verify(http.HandlerFunc(rakerServer.HistoryPage)))
	mux.Handle("/instagram", rakerServer.Verify(http.HandlerFunc(rakerServer.InstagramPage)))
	mux.Handle("/highlight", rakerServer.Verify(http.HandlerFunc(rakerServer.HighlightPage)))
	mux.Handle("/story", rakerServer.Verify(http.HandlerFunc(rakerServer.StoryPage)))
	mux.Handle("/tiktok", rakerServer.Verify(http.HandlerFunc(rakerServer.TikTokPage)))
	mux.Handle("/vsco", rakerServer.Verify(http.HandlerFunc(rakerServer.VSCOPage)))

	rakerServer.HTTPServer = http.Server{
		Addr:    fmt.Sprintf(":%d", configuration.Port),
		Handler: NewLoggerMiddleware(mux),
		ErrorLog: log.Default().StandardLog(log.StandardLogOptions{
			ForceLevel: log.ErrorLevel,
		}),
	}

	return &rakerServer, nil
}

var configuration = Configuration{
	URI:         "mongodb://localhost:27017",
	Database:    "raker",
	Storage:     ".",
	Directories: false,
	Port:        4100,
}

func init() {
	// viper.SetEnvPrefix("raker")
	viper.AutomaticEnv()
	viper.BindEnv("SECRET")
	viper.BindEnv("URI")
	viper.BindEnv("DATABASE")
	viper.BindEnv("STORAGE")
	viper.BindEnv("DIRECTORIES")
	viper.BindEnv("PORT")
	viper.SetConfigName(".raker")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
}
