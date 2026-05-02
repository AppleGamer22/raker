package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	rpcPath, handler := v1connect.NewRakerServerHandler(
		&rakerServer,
		connect.WithInterceptors(rakerServer.NewAuthInterceptor(), validate.NewInterceptor()),
	)

	mux := http.NewServeMux()
	// Connect RPC
	mux.Handle(fmt.Sprintf("/api%s", rpcPath), http.StripPrefix("/api", handler))
	// Storage
	mux.Handle("/api/storage/", http.StripPrefix("/api/storage", rakerServer.NewStorageHandler(rakerServer.Configuration.Storage, rakerServer.Configuration.Directories)))
	// React client: serve static files from dist, but fall back to index.html
	fileServer := http.FileServer(http.Dir("dist"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the requested file; if it doesn't exist, serve index.html
		reqPath := strings.TrimPrefix(r.URL.Path, "/")
		fullPath := filepath.Join("dist", reqPath)
		if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join("dist", "index.html"))
	})

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
