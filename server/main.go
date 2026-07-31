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
	"time"

	"github.com/AppleGamer22/raker/server/handlers"
	"github.com/AppleGamer22/raker/shared"
	"github.com/charmbracelet/log"
	"go.yaml.in/yaml/v3"
)

func main() {
	rakerServer, err := handlers.NewRakerServer()
	if err != nil {
		log.Fatal(err)
	}
	defer rakerServer.DBClient.Close()

	log.Infof("raker %s %s (%s/%s)", shared.Version, shared.Hash, runtime.GOOS, runtime.GOARCH)
	log.Infof("Storage path: %s", rakerServer.Storage)
	if rakerServer.Directories {
		log.Info("allowing directory listing")
	}
	log.Infof("database URI: %s", rakerServer.URI)
	log.Infof("Server is listening at http://localhost:%d", rakerServer.Port)
	log.Infof("WebAuthn RPID: %s", rakerServer.RPID)
	if len(rakerServer.RPOrigins) > 0 {
		yamlRPOrigins, _ := yaml.Marshal(rakerServer.RPOrigins)
		log.Infof("WebAuthn RPOrigins: %s", string(yamlRPOrigins))
	}

	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		if err := rakerServer.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error(err)
			signals <- os.Interrupt
		}
	}()

	<-signals
	fmt.Print("\r")
	log.Warn("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rakerServer.HTTPServer.Shutdown(ctx); err != nil {
		log.Warn(err)
	}
}
