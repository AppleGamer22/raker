package conf

import (
	"fmt"
	"io"
	"net/http"
	"os"

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

	_, err = io.Copy(file, response.Body)
	return err
}
