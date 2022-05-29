package handlers

import (
	"log"
	"net/http"
	"os"
	"path"
)

type storageHandler struct {
	root        string
	directories bool
	fileServer  http.Handler
}

func NewStorageHandler(root string, directories bool) storageHandler {
	return storageHandler{
		root:        root,
		directories: directories,
		fileServer:  http.FileServer(http.Dir(root)),
	}
}

func (handler storageHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	mediaPath := path.Join(handler.root, request.URL.Path)
	info, err := os.Stat(mediaPath)
	switch request.Method {
	case http.MethodGet:
		if handler.directories || (err == nil && !info.IsDir()) {
			handler.fileServer.ServeHTTP(writer, request)
		} else {
			if err != nil {
				log.Println(err, request.URL.Path)
			}
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}

}
