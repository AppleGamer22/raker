package handlers

import (
	"log"
	"net/http"
)

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
)

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

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		statusWriter := newStatusLogWriter(writer)
		handler.ServeHTTP(statusWriter, request)
		switch statusWriter.Status / 100 {
		case 4, 5:
			log.Println(colorRed, statusWriter.Status, colorReset, request.Method, request.RequestURI, request.RemoteAddr)
		default:
			log.Println(colorGreen, statusWriter.Status, colorReset, request.Method, request.RequestURI, request.RemoteAddr)
		}
	})
}
