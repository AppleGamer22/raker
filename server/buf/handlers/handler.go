package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/AppleGamer22/raker/server/authenticator"
	"github.com/AppleGamer22/raker/server/buf/connect/raker/v1/v1connect"
	v1 "github.com/AppleGamer22/raker/server/buf/proto/raker/v1"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/emptypb"

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

// EditCategory implements [v1connect.RakerServerHandler].
func (r *RakerServer) EditCategory(context.Context, *v1.EditCategoryRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// RemoveFile implements [v1connect.RakerServerHandler].
func (r *RakerServer) RemoveFile(context.Context, *v1.RemoveFileRequest) (*v1.ScrapeResponse, error) {
	panic("unimplemented")
}

// ScrapeHighlight implements [v1connect.RakerServerHandler].
func (r *RakerServer) ScrapeHighlight(context.Context, *v1.UnaryScrapeRequest) (*v1.ScrapeResponse, error) {
	panic("unimplemented")
}

// ScrapeInstagram implements [v1connect.RakerServerHandler].
func (r *RakerServer) ScrapeInstagram(context.Context, *v1.UnaryScrapeRequest) (*v1.ScrapeResponse, error) {
	panic("unimplemented")
}

// ScrapeSnapchat implements [v1connect.RakerServerHandler].
func (r *RakerServer) ScrapeSnapchat(context.Context, *v1.UnaryScrapeRequest) (*v1.ScrapeResponse, error) {
	panic("unimplemented")
}

// ScrapeStory implements [v1connect.RakerServerHandler].
func (r *RakerServer) ScrapeStory(context.Context, *v1.UnaryScrapeRequest) (*v1.ScrapeResponse, error) {
	panic("unimplemented")
}

// ScrapeTikTok implements [v1connect.RakerServerHandler].
func (r *RakerServer) ScrapeTikTok(context.Context, *v1.BinaryScrapeRequest) (*v1.ScrapeResponse, error) {
	panic("unimplemented")
}

// ScrapeVSCO implements [v1connect.RakerServerHandler].
func (r *RakerServer) ScrapeVSCO(context.Context, *v1.BinaryScrapeRequest) (*v1.ScrapeResponse, error) {
	panic("unimplemented")
}

// SearchHistory implements [v1connect.RakerServerHandler].
func (r *RakerServer) SearchHistory(context.Context, *v1.HistoryRequest) (*v1.HistoryResponse, error) {
	panic("unimplemented")
}

// SearchHistoryOwners implements [v1connect.RakerServerHandler].
func (r *RakerServer) SearchHistoryOwners(context.Context, *v1.HistoryRequest) (*v1.HistoryOwnersResponse, error) {
	panic("unimplemented")
}

// SignInInstagram implements [v1connect.RakerServerHandler].
func (r *RakerServer) SignInInstagram(context.Context, *v1.SignUpRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// SignUpInstagram implements [v1connect.RakerServerHandler].
func (r *RakerServer) SignUpInstagram(context.Context, *v1.SignUpRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// UpdateCategories implements [v1connect.RakerServerHandler].
func (r *RakerServer) UpdateCategories(context.Context, *v1.UpdateCategoriesRequest) (*v1.ScrapeResponse, error) {
	panic("unimplemented")
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

	mux := http.NewServeMux()

	path, handler := v1connect.NewRakerServerHandler(
		&rakerServer,
		connect.WithInterceptors(validate.NewInterceptor()),
	)

	mux.Handle(path, handler)
	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)

	rakerServer.HTTPServer = http.Server{
		Addr:      fmt.Sprintf(":%d", rakerServer.Configuration.Port),
		Handler:   mux,
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
