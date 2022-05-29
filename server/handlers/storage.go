package handlers

import (
	"net/http"
	"strings"
)

type storageHandler struct {
	prefix      string
	directories bool
	fileServer  http.Handler
}

func NewStorageHandler(prefix, root string, directories bool) storageHandler {
	return storageHandler{
		prefix:      prefix,
		directories: directories,
		fileServer:  http.FileServer(http.Dir(root)),
	}
}

func (handler storageHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	urlString := request.URL.String()
	validURL := strings.HasSuffix(urlString, ".jpg") || strings.HasSuffix(urlString, ".mp4") || strings.HasSuffix(urlString, ".webp") || strings.HasSuffix(urlString, ".webm")
	switch request.Method {
	case http.MethodGet:
		if handler.directories || validURL {
			http.StripPrefix(handler.prefix, handler.fileServer).ServeHTTP(writer, request)
		} else {
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}

}
