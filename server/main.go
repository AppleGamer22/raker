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

	"github.com/AppleGamer22/raker/server/configuration"
	"github.com/AppleGamer22/raker/shared"
	"github.com/charmbracelet/log"
)

func main() {
	rakerServer, err := configuration.NewRakerServer()
	if err != nil {
		log.Fatal(err)
	}
	defer rakerServer.DBClient.Disconnect(context.Background())

	log.Infof("raker %s %s (%s/%s)", shared.Version, shared.Hash, runtime.GOOS, runtime.GOARCH)
	log.Infof("Storage path: %s", rakerServer.Storage)
	if rakerServer.Directories {
		log.Info("allowing directory listing")
	}
	log.Infof("MongoDB database URI: %s", rakerServer.URI)
	log.Infof("MongoDB database: %s", rakerServer.Database)
	log.Infof("Server is listening at http://localhost:%d", rakerServer.Port)

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
