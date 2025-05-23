package configuration

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/AppleGamer22/raker/server/cleaner"
	db "github.com/AppleGamer22/raker/server/db/mongo"
	"github.com/AppleGamer22/raker/shared"
	"github.com/AppleGamer22/raker/shared/types"
	"github.com/bep/imagemeta"
	"github.com/charmbracelet/log"
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
			log.Error(err, "URL", escapedURL)
		}
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func (handler *storageHandler) Save(user db.User, media, owner, fileName, URL string, cookies []*http.Cookie) error {
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
	response, err := http.DefaultClient.Do(request)
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

func (handler *storageHandler) SaveBundle(user db.User, media, owner string, fileNames, URLs []string, cookies []*http.Cookie) ([]string, []error) {
	if len(URLs) != len(fileNames) {
		return []string{}, []error{errors.New("unequal length URLs & file names slices")}
	}

	count := len(URLs)
	var wg sync.WaitGroup
	wg.Add(count)
	var mutex sync.Mutex
	errs := make([]error, 0, count)

	for i := 0; i < count; i++ {
		URL := URLs[i]
		fileName := fileNames[i]
		go func(fileName, URL string, i int) {
			if err := handler.Save(user, media, owner, fileName, URL, cookies); err != nil {
				mutex.Lock()
				errs = append(errs, err)
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

	return sucessfulFileNames, errs
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
	filePath := path.Join(user.ID.Hex(), media, owner, fileName)
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

	if err := imagemeta.Decode(imagemeta.Options{R: file, HandleTag: handleTag, ImageFormat: imagemeta.JPEG}); err != nil {
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
