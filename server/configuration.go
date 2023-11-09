package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AppleGamer22/raker/server/authenticator"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/AppleGamer22/raker/server/handlers"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

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
