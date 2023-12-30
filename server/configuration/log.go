package configuration

import (
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

type LoggerMiddleware struct {
	handler http.Handler
}

func NewLoggerMiddleware(handler http.Handler) *LoggerMiddleware {
	return &LoggerMiddleware{handler: handler}
}

func (lm *LoggerMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	statusWriter := newStatusLogWriter(writer)
	startTime := time.Now()
	lm.handler.ServeHTTP(statusWriter, request)
	duration := time.Since(startTime).Round(time.Microsecond)
	switch statusWriter.Status / 100 {
	case 4, 5:
		log.Errorf("%d %v %s %s %s", statusWriter.Status, duration, request.Method, request.URL.Path, request.RemoteAddr)
	default:
		log.Infof("%d %v %s %s %s", statusWriter.Status, duration, request.Method, request.URL.Path, request.RemoteAddr)
	}
}

// StatusWriter type adds a Status property to Go's http.ResponseWriter type.
type statusLogWriter struct {
	http.ResponseWriter
	Status int
}

func newStatusLogWriter(writer http.ResponseWriter) *statusLogWriter {
	return &statusLogWriter{
		ResponseWriter: writer,
		Status:         http.StatusOK,
	}
}

// WriteHeader overrides the http.ResponseWriter's WriteHeader method
func (writer *statusLogWriter) WriteHeader(status int) {
	writer.ResponseWriter.WriteHeader(status)
	writer.Status = status
}
