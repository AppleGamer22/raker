package conf

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/AppleGamer22/rake/shared"
	"github.com/AppleGamer22/rake/shared/types"
)

func Save(media, fileName, URL string) error {
	if !types.ValidMediaType(media) {
		return fmt.Errorf("invalid media type: %s", media)
	}

	_, err := os.Stat(fileName)
	if err == nil {
		return fmt.Errorf("file %s already exists", fileName)
	}

	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return err
	}

	if media == types.TikTok {
		request.Header.Add("Range", "bytes=0-")
		sessionCookie := http.Cookie{
			Name:     "sessionid",
			Value:    Configuration.TikTok,
			Domain:   ".tiktok.com",
			HttpOnly: true,
		}
		request.AddCookie(&sessionCookie)
	}
	request.Header.Add("User-Agent", shared.UserAgent)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, response.Body); err != nil {
		return err
	}

	log.Printf("saved %s at the current directory\n", fileName)
	return nil
}

func SaveBundle(media string, fileNames, URLs []string) []error {
	if len(URLs) != len(fileNames) {
		return []error{errors.New("unequal length URLs & file names slices")}
	}

	count := len(URLs)
	var wg sync.WaitGroup
	wg.Add(count)
	var mutex sync.Mutex
	errs := make([]error, 0, count)

	for i := 0; i < count; i++ {
		URL := URLs[i]
		fileName := fileNames[i]
		go func() {
			if err := Save(media, fileName, URL); err != nil {
				mutex.Lock()
				errs = append(errs, err)
				mutex.Unlock()
			}
			wg.Done()
		}()
	}

	wg.Wait()
	return errs
}
