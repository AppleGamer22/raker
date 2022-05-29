package handlers

import (
	"log"
	"net/http"
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handler.ServeHTTP(writer, request)
		log.Println(request.Method, request.RequestURI, request.RemoteAddr)
	})
}
