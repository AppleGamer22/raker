package configuration

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AppleGamer22/raker/server/authenticator"
	db "github.com/AppleGamer22/raker/server/db/mongo"
	"github.com/AppleGamer22/raker/shared"
	"github.com/AppleGamer22/raker/shared/types"

	"github.com/charmbracelet/log"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type contextKey int

const authenticatedUserKey contextKey = 0

type Configuration struct {
	Secret       string
	URI          string
	Database     string
	Username     string
	Password     string
	Storage      string
	Directories  bool
	SecureCookie bool
	Port         uint
}

type RakerServer struct {
	Configuration
	DBClient      *mongo.Client
	Users         *mongo.Collection
	Histories     *mongo.Collection
	Authenticator authenticator.Authenticator
	WebAuthn      *webauthn.WebAuthn
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

	if rakerServer.Configuration.Secret == "" && !viper.IsSet("secret") {
		log.Fatal("A JWT secret must be set via a config file or an environment variable")
	}

	rakerServer.Authenticator = authenticator.New(rakerServer.Configuration.Secret)

	dbClient, database, err := db.Connect(
		rakerServer.Configuration.URI,
		rakerServer.Configuration.Database,
		rakerServer.Configuration.Username,
		rakerServer.Configuration.Password,
	)
	if err != nil {
		log.Fatal(err)
	}
	rakerServer.DBClient = dbClient
	// remember to defer client.Disconnet
	rakerServer.Histories = database.Collection("histories")
	rakerServer.Users = database.Collection("users")

	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth/sign_up/find/instagram", rakerServer.InstagramSignUp)
	mux.HandleFunc("/api/auth/sign_in/find/instagram", rakerServer.InstagramSignIn)
	mux.Handle("/api/auth/update/find/instagram", rakerServer.Verify(true, http.HandlerFunc(rakerServer.InstagramUpdateCredentials)))
	mux.Handle("/api/auth/sign_out/find/instagram", rakerServer.Verify(true, http.HandlerFunc(rakerServer.InstagramSignOut)))
	mux.Handle("/api/categories", rakerServer.Verify(true, http.HandlerFunc(rakerServer.Categories)))
	mux.Handle("/api/history", rakerServer.Verify(true, http.HandlerFunc(rakerServer.History)))
	// mux.HandleFunc("/api/info", rakerServer.Information)
	mux.Handle("/api/storage/", http.StripPrefix("/api/storage", rakerServer.Verify(true, NewStorageHandler(rakerServer.Configuration.Storage, rakerServer.Configuration.Directories))))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.Handle("/favicon.ico", http.RedirectHandler("/assets/icons/favicon.ico", http.StatusPermanentRedirect))
	mux.Handle("/robots.txt", http.RedirectHandler("/assets/robots.txt", http.StatusPermanentRedirect))

	// TODO: look into https://htmx.org/attributes/hx-indicator/

	mux.Handle("/", rakerServer.Verify(false, http.HandlerFunc(rakerServer.AuthenticationPage)))
	mux.Handle("/history", rakerServer.Verify(true, http.HandlerFunc(rakerServer.HistoryPage)))
	mux.Handle("GET /find/{type}", rakerServer.Verify(true, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var (
			user    db.User
			history db.History
			errs    []error
		)

		switch request.PathValue("type") {
		case types.Instagram:
			user, history, errs = rakerServer.instagram(request)
		case types.Highlight:
			user, history, errs = rakerServer.highlight(request)
		case types.Story:
			user, history, errs = rakerServer.story(request)
		case types.TikTok:
			user, history, errs = rakerServer.tiktok(request)
		case types.VSCO:
			user, history, errs = rakerServer.vsco(request)
		default:
			http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
			return
		}

		if len(errs) > 0 {
			writer.WriteHeader(http.StatusBadRequest)
			for _, err := range errs {
				log.Error(err)
			}
		}

		historyDisplay := db.HistoryDisplay{
			History:            history,
			Errors:             errs,
			Version:            shared.Version,
			SelectedCategories: user.SelectedCategories(history.Categories),
		}

		if err := templates.ExecuteTemplate(writer, "history.html", historyDisplay); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
	})))

	mux.Handle("GET /find/{type}/htmx", rakerServer.Verify(true, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var (
			user    db.User
			history db.History
			errs    []error
		)

		switch request.PathValue("type") {
		case types.Instagram:
			user, history, errs = rakerServer.instagram(request)
		case types.Highlight:
			user, history, errs = rakerServer.highlight(request)
		case types.Story:
			user, history, errs = rakerServer.story(request)
		case types.TikTok:
			user, history, errs = rakerServer.tiktok(request)
		case types.VSCO:
			user, history, errs = rakerServer.vsco(request)
		default:
			http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
			return
		}

		if len(errs) > 0 {
			// writer.WriteHeader(http.StatusBadRequest)
			for _, err := range errs {
				log.Error(err)
			}
		}

		historyDisplay := db.HistoryDisplay{
			History:            history,
			Errors:             errs,
			SelectedCategories: user.SelectedCategories(history.Categories),
		}

		if err := templates.ExecuteTemplate(writer, "history_result.html", historyDisplay); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			log.Error(err)
		}
	})))

	rakerServer.HTTPServer = http.Server{
		Addr:    fmt.Sprintf(":%d", rakerServer.Configuration.Port),
		Handler: NewLoggerMiddleware(mux),
		ErrorLog: log.Default().StandardLog(log.StandardLogOptions{
			ForceLevel: log.ErrorLevel,
		}),
	}

	return &rakerServer, nil
}

func init() {
	// viper.SetEnvPrefix("raker")
	viper.AutomaticEnv()
	viper.BindEnv("SECRET")
	viper.BindEnv("URI")
	viper.BindEnv("DATABASE")
	viper.BindEnv("STORAGE")
	viper.BindEnv("DIRECTORIES")
	viper.BindEnv("SECURE_COOKIES")
	viper.BindEnv("PORT")
	viper.SetConfigName(".raker")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
}
