package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/AppleGamer22/rake/server/cleaner"
	"github.com/AppleGamer22/rake/server/db"
	"github.com/AppleGamer22/rake/shared"
	"github.com/AppleGamer22/rake/shared/types"
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
	mediaPath = cleaner.Path(mediaPath)

	info, err := os.Stat(mediaPath)
	if handler.directories || (err == nil && !info.IsDir()) {
		handler.fileServer.ServeHTTP(writer, request)
	} else {
		if err != nil {
			escapedURL := cleaner.Line(request.URL.Path)
			log.Println(err, escapedURL)
		}
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func (handler *storageHandler) Save(user db.User, media, owner, fileName, URL string) error {
	if !types.ValidMediaType(media) {
		return fmt.Errorf("invalid media type: %s", media)
	}

	filePath := path.Join(user.ID.Hex(), media, owner, fileName)
	mediaPath := path.Join(handler.root, filePath)
	mediaPath = cleaner.Path(mediaPath)

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

	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return err
	}

	if media == types.TikTok {
		request.Header.Add("Range", "bytes=0-")
		if user.TikTok != "" {
			sessionCookie := http.Cookie{
				Name:     "sessionid",
				Value:    user.TikTok,
				Domain:   ".tiktok.com",
				HttpOnly: true,
			}
			request.AddCookie(&sessionCookie)
		}
	}
	request.Header.Add("User-Agent", shared.UserAgent)
	response, err := http.DefaultClient.Do(request)
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

func (handler *storageHandler) Delete(user db.User, media, owner, fileName string) error {
	if !types.ValidMediaType(media) {
		return fmt.Errorf("invalid media type: %s", media)
	}

	filePath := path.Join(user.ID.Hex(), media, owner, fileName)
	mediaPath := path.Join(handler.root, filePath)
	mediaPath = cleaner.Path(mediaPath)

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
