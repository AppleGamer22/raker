package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/AppleGamer22/rake/server/db"
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
	switch request.Method {
	case http.MethodGet:
		mediaPath := path.Join(handler.root, request.URL.Path)
		info, err := os.Stat(mediaPath)
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

func (handler storageHandler) Save(media, owner, fileName, URL string) error {
	if !db.ValidMediaType(media) {
		return fmt.Errorf("invalid media type: %s", media)
	}

	filePath := path.Join(media, owner, fileName)
	mediaPath := path.Join(handler.root, filePath)

	_, err := os.Stat(mediaPath)
	if err == nil {
		return fmt.Errorf("file %s already exists", filePath)
	}

	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(mediaPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	return err
}

func (handler storageHandler) Delete(media, owner, fileName, URL string) error {
	if !db.ValidMediaType(media) {
		return fmt.Errorf("invalid media type: %s", media)
	}

	filePath := path.Join(media, owner, fileName)
	mediaPath := path.Join(handler.root, filePath)

	_, err := os.Stat(mediaPath)
	if err != nil {
		return fmt.Errorf("file %s does not exists", filePath)
	}

	err = os.Remove(mediaPath)

	directoryName := path.Dir(mediaPath)
	files, err2 := os.ReadDir(directoryName)
	if err2 != nil {
		err = fmt.Errorf("%v\n%v", err, err2)
	}
	if len(files) == 0 {
		if err3 := os.Remove(directoryName); err3 != nil {
			err = fmt.Errorf("%v\n%v", err, err3)
		}
	}

	return err
}
