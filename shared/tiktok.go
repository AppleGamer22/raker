package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type TikTokPost struct {
	ItemModule map[string]struct {
		Author string `json:"author"`
		Video  struct {
			DownloadAddress string `json:"downloadAddr"`
		} `json:"video"`
	}
}

func TikTok(owner, post string) (URL string, username string, err error) {
	postURL := fmt.Sprintf("https://www.tiktok.com/@%s/video/%s", owner, post)
	request, err := http.NewRequest(http.MethodGet, postURL, nil)
	if err != nil {
		return URL, username, err
	}

	request.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return URL, username, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return URL, username, err
	}

	re := regexp.MustCompile(`<script id=\"SIGI_STATE\" type=\"application/json\">(.*?)</script>`)
	script := re.FindString(string(body))
	if script == "" {
		return URL, username, errors.New("could not find JSON")
	}

	jsonText := script[len(`<script id="SIGI_STATE" type="application/json">`) : len(script)-len("</script>")]
	var tiktokPost TikTokPost
	if err := json.Unmarshal([]byte(jsonText), &tiktokPost); err != nil {
		return URL, username, err
	}

	item := tiktokPost.ItemModule[post]
	username = item.Author
	URL = item.Video.DownloadAddress

	return URL, username, err
}
