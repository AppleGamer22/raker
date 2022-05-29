package handlers

import (
	"net/http"
	"strings"
)

type storageServer struct {
	prefix     string
	fileServer http.Handler
}

func NewStorageServer(prefix, root string) storageServer {
	return storageServer{
		prefix:     prefix,
		fileServer: http.FileServer(http.Dir(root)),
	}
}

func (server storageServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	urlString := request.URL.String()
	validURL := strings.HasSuffix(urlString, ".jpg") || strings.HasSuffix(urlString, ".mp4") || strings.HasSuffix(urlString, ".webp") || strings.HasSuffix(urlString, ".webm")
	if request.Method == "GET" && validURL {
		http.StripPrefix(server.prefix, server.fileServer).ServeHTTP(writer, request)
	} else {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}
