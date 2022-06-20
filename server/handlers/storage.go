package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"

	"github.com/AppleGamer22/rake/server/db"
)

type storageHandler struct {
	root        string
	directories bool
	fileServer  http.Handler
}

var StorageHandler storageHandler

func NewStorageHandler(root string, directories bool) storageHandler {
	StorageHandler = storageHandler{
		root:        root,
		directories: directories,
		fileServer:  http.FileServer(http.Dir(root)),
	}
	return StorageHandler
}

func (handler storageHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	mediaPath := path.Join(handler.root, request.URL.Path)

	if runtime.GOOS == "windows" {
		mediaPath = strings.ReplaceAll(mediaPath, `..\`, "")
		mediaPath = regexp.MustCompile(`[A-Z]:`).ReplaceAllString(mediaPath, "")
	} else {
		mediaPath = strings.ReplaceAll(mediaPath, "../", "")
	}

	info, err := os.Stat(mediaPath)
	switch request.Method {
	case http.MethodDelete:
		pathComponents := strings.Split(request.URL.Path, "/")
		if len(pathComponents) != 3 {
			http.Error(writer, "invalid path", http.StatusBadRequest)
			return
		}
		if err := handler.Delete(pathComponents[0], pathComponents[1], pathComponents[2]); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		}
	case http.MethodPatch:
	default:
		if handler.directories || (err == nil && !info.IsDir()) {
			handler.fileServer.ServeHTTP(writer, request)
		} else {
			if err != nil {
				escapedURL := strings.Replace(request.URL.Path, "\n", "", -1)
				escapedURL = strings.Replace(escapedURL, "\r", "", -1)
				log.Println(err, escapedURL)
			}
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
}

func (handler *storageHandler) Save(media, owner, fileName, URL string) error {
	if !db.ValidMediaType(media) {
		return fmt.Errorf("invalid media type: %s", media)
	}

	filePath := path.Join(media, owner, fileName)
	mediaPath := path.Join(handler.root, filePath)

	_, err := os.Stat(mediaPath)
	if err == nil {
		return fmt.Errorf("file %s already exists", filePath)
	}

	directoryName := path.Dir(mediaPath)
	if _, err := os.Stat(directoryName); err != nil {
		const userGroupReadable = 660
		if err := os.MkdirAll(directoryName, userGroupReadable); err != nil {
			return err
		}
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

func (handler *storageHandler) Delete(media, owner, fileName string) error {
	if !db.ValidMediaType(media) {
		return fmt.Errorf("invalid media type: %s", media)
	}

	filePath := path.Join(media, owner, fileName)
	mediaPath := path.Join(handler.root, filePath)

	_, err := os.Stat(mediaPath)
	if err != nil {
		return fmt.Errorf("file %s does not exists", filePath)
	}

	if err := os.Remove(mediaPath); err != nil {
		return err
	}

	directoryName := path.Dir(mediaPath)
	files, err := os.ReadDir(directoryName)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		if err := os.Remove(directoryName); err != nil {
			return err
		}
	}

	return nil
}
