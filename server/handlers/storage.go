package handlers

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/AppleGamer22/raker/server/cleaner"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/AppleGamer22/raker/shared"
	"github.com/AppleGamer22/raker/shared/types"
	"github.com/bep/imagemeta"
	"github.com/charmbracelet/log"

	"context"

	v1 "github.com/AppleGamer22/raker/server/buf/proto/raker/v1"
)

// RemoveFile implements [v1connect.RakerServerHandler].
func (r *RakerServer) RemoveFile(context.Context, *v1.RemoveFileRequest) (*v1.ScrapeResponse, error) {
	panic("unimplemented")
}

type storageHandler struct {
	root        string
	directories bool
	fileServer  http.Handler
	server      *RakerServer
}

var StorageHandler storageHandler

func (server *RakerServer) NewStorageHandler(root string, directories bool) storageHandler {
	StorageHandler = storageHandler{
		root:        root,
		directories: directories,
		fileServer:  http.FileServer(http.Dir(root)),
		server:      server,
	}
	return StorageHandler
}

func (handler storageHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("jwt")
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	u, err := handler.server.GetUserFromCookie(cookie)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if !strings.HasPrefix(request.URL.Path, fmt.Sprintf("/%s/", u.Username)) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	mediaPath := path.Join(handler.root, request.URL.Path)
	mediaPath = cleaner.Path(mediaPath)

	info, err := os.Stat(mediaPath)
	if handler.directories || (err == nil && !info.IsDir()) {
		handler.fileServer.ServeHTTP(writer, request)
	} else {
		if err != nil {
			escapedURL := cleaner.Line(request.URL.Path)
			log.Error(err, "URL", escapedURL)
		}
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func (handler *storageHandler) Save(user db.User, media db.PostType, owner, fileName, URL string, cookies []*http.Cookie) error {
	filePath := path.Join(user.Username, string(media), owner, fileName)
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

	switch media {
	case types.TikTok:
		request.Header.Add("Range", "bytes=0-")
		for _, cookie := range cookies {
			request.AddCookie(cookie)
		}
		request.Header.Add("referer", "https://www.tiktok.com/")
	case types.VSCO:
		request.Header.Add("referer", "https://vsco.co/")
	}

	request.Header.Add("User-Agent", shared.UserAgent)

	client := http.DefaultClient
	if media == types.VSCO {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS13,
				},
			},
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	statusClass := response.StatusCode / 100
	if statusClass == 4 || statusClass == 5 {
		return fmt.Errorf("response of %d instead of media", response.StatusCode)
	}

	if media == types.VSCO && response.Header.Get("Content-Type") == "video/MP2T" {
		log.Debugf("decoding stream from %s", media)
		return shared.Stream2MP4(response.Body, mediaPath)
	}

	file, err := os.Create(mediaPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, response.Body); err != nil {
		return err
	}

	log.Debugf("saved %s", filePath)
	return err
}

func (handler *storageHandler) SaveBundle(user db.User, media db.PostType, owner string, fileNames, URLs []string, cookies []*http.Cookie) ([]string, error) {
	if len(URLs) != len(fileNames) {
		return []string{}, errors.New("unequal length URLs & file names slices")
	}

	count := len(URLs)
	var wg sync.WaitGroup
	wg.Add(count)
	var mutex sync.Mutex
	var err error = nil

	for i := 0; i < count; i++ {
		URL := URLs[i]
		fileName := fileNames[i]
		go func(fileName, URL string, i int) {
			if err2 := handler.Save(user, media, owner, fileName, URL, cookies); err != nil {
				mutex.Lock()
				err = errors.Join(err, err2)
				fileNames[i] = ""
				mutex.Unlock()
			}
			wg.Done()
		}(fileName, URL, i)
	}

	wg.Wait()

	sucessfulFileNames := make([]string, 0, count)
	for _, fileName := range fileNames {
		if fileName != "" {
			sucessfulFileNames = append(sucessfulFileNames, fileName)
		}
	}

	return sucessfulFileNames, err
}

func (handler *storageHandler) Delete(user db.User, media, owner, fileName string) error {
	if !types.ValidMediaType(media) {
		return fmt.Errorf("invalid media type: %s", media)
	}

	filePath := path.Join(user.Username, media, owner, fileName)
	mediaPath := path.Join(handler.root, filePath)
	mediaPath = cleaner.Path(mediaPath)

	_, err := os.Stat(mediaPath)
	if err != nil {
		return fmt.Errorf("file %s does not exists", filePath)
	}

	if err := os.Remove(mediaPath); err != nil {
		return err
	}
	log.Warnf("deleted %s", filePath)

	directoryName := path.Dir(mediaPath)
	files, err := os.ReadDir(directoryName)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		if err := os.Remove(directoryName); err != nil {
			return err
		}
		log.Warnf("deleted %s", directoryName)
	}

	return nil
}

func (handler *storageHandler) LocationEXIF(user db.User, media, owner, fileName string) (float64, float64) {
	filePath := path.Join(user.Username, media, owner, fileName)
	mediaPath := path.Join(handler.root, filePath)
	mediaPath = cleaner.Path(mediaPath)
	file, err := os.Open(mediaPath)
	if err != nil {
		log.Error(err)
		return 0, 0
	}
	defer file.Close()

	var tags imagemeta.Tags
	handleTag := func(ti imagemeta.TagInfo) error {
		tags.Add(ti)
		return nil
	}

	if _, err := imagemeta.Decode(imagemeta.Options{R: file, HandleTag: handleTag, ImageFormat: imagemeta.JPEG}); err != nil {
		log.Error(err)
		return 0, 0
	}

	latitude, longitude, err := tags.GetLatLong()
	if err != nil {
		log.Error(err)
		return 0, 0
	}

	return latitude, longitude
}
