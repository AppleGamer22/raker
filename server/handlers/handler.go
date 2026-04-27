package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/AppleGamer22/raker/server/authenticator"
	"github.com/AppleGamer22/raker/server/buf/connect/raker/v1/v1connect"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/spf13/viper"

	"github.com/charmbracelet/log"
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
	DBConnection  *sql.DB
	DBClient      *db.Queries
	Authenticator authenticator.Authenticator
	WebAuthn      *webauthn.WebAuthn
	HTTPServer    http.Server
}

// type rakerServerHandlerAdapter struct {
// 	*RakerServer
// }

// func (a *rakerServerHandlerAdapter) SignInInstagram(ctx context.Context, req *v1.SignInRequest) (*emptypb.Empty, error) {
// 	resp, err := a.RakerServer.SignInInstagram(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if resp == nil || resp.Msg == nil {
// 		return &emptypb.Empty{}, nil
// 	}
// 	return resp.Msg, nil
// }

func NewRakerServer() (*RakerServer, error) {
	rakerServer := RakerServer{
		Configuration: Configuration{
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

	connection, err := sql.Open("postgres", rakerServer.URI)

	if err != nil {
		log.Fatal(err)
	}

	if err := connection.Ping(); err != nil {
		log.Fatal(err)
	}
	rakerServer.DBConnection = connection
	pgdb := db.New(connection)
	rakerServer.DBClient = pgdb

	path, handler := v1connect.NewRakerServerHandler(
		&rakerServer,
		connect.WithInterceptors(rakerServer.NewAuthInterceptor(), validate.NewInterceptor()),
	)

	mux := http.NewServeMux()
	mux.Handle(fmt.Sprintf("/api%s", path), http.StripPrefix("/api", handler))
	mux.Handle("/api/storage/", http.StripPrefix("/api/storage", rakerServer.NewStorageHandler(rakerServer.Configuration.Storage, rakerServer.Configuration.Directories)))

	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)

	rakerServer.HTTPServer = http.Server{
		Addr:      fmt.Sprintf(":%d", rakerServer.Configuration.Port),
		Handler:   NewLoggerMiddleware(mux),
		Protocols: protocols,
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
