package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/AppleGamer22/raker/server/authenticator"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/AppleGamer22/raker/server/handlers"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
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

func (rs *RakerServer) Verify(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		jwtCookie, err := request.Cookie("jwt")
		if err != nil {
			http.Error(writer, "credential update failed", http.StatusUnauthorized)
			log.Error(err)
			return
		}

		U_ID, username, err := rs.Authenticator.Parse(jwtCookie.Value)
		if err != nil {
			http.Error(writer, "credential update failed", http.StatusUnauthorized)
			log.Error(err)
			return
		}

		filter := bson.M{
			"_id":      U_ID,
			"username": username,
		}
		var user db.User
		if err := rs.Users.FindOne(context.Background(), filter).Decode(&user); err != nil {
			http.Error(writer, "credential update failed", http.StatusUnauthorized)
			log.Error(err)
			return
		}
		// https://drstearns.github.io/tutorials/gomiddleware/#secmiddlewareandrequestscopedvalues
		ctxWithUser := context.WithValue(request.Context(), authenticatedUserKey, user)
		requestWithUser := request.WithContext(ctxWithUser)
		handler.ServeHTTP(writer, requestWithUser)
	})
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

	dbClient, err := db.Connect(configuration.URI, configuration.Database)
	if err != nil {
		log.Fatal(err)
	}
	// remember to defer client.Close()
	rakerServer.DBClient = dbClient

	mux := http.NewServeMux()

	rakerServer.HTTPServer = http.Server{
		Addr:    fmt.Sprintf(":%d", configuration.Port),
		Handler: handlers.NewLoggerMiddleware(mux),
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
