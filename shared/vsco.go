package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type VSCOPost struct {
	Medias struct {
		ByID map[string]struct {
			Media struct {
				PermaSubdomain string `json:"permaSubdomain"`
				ResponsiveURL  string `json:"responsiveUrl"`
				VideoURL       string `json:"videoUrl"`
			} `json:"media"`
		} `json:"byId"`
	} `json:"medias"`
}

var vsco_regexp = regexp.MustCompile(`<script>window\.__PRELOADED_STATE__ =(.*?)</script>`)

func VSCO(owner, post string) (URL string, username string, err error) {
	postURL := fmt.Sprintf("https://vsco.co/%s/media/%s", owner, post)
	response, err := http.Get(postURL)
	if err != nil {
		return URL, username, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return URL, username, err
	}

	script := vsco_regexp.FindString(string(body))
	if script == "" {
		return URL, username, errors.New("could not find JSON")
	}

	jsonText := script[len("<script>window.__PRELOADED_STATE__ =") : len(script)-len("</script>")]
	var vscoPost VSCOPost
	if err := json.Unmarshal([]byte(jsonText), &vscoPost); err != nil {
		return URL, username, err
	}

	media := vscoPost.Medias.ByID[post]
	username = media.Media.PermaSubdomain

	if len(media.Media.VideoURL) > 0 {
		URL = fmt.Sprintf("https://%s", media.Media.VideoURL)
	} else {
		URL = fmt.Sprintf("https://%s", media.Media.ResponsiveURL)
	}

	return URL, username, err
}
