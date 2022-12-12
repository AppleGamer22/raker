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

type TikTok struct {
	SessionID  string
	ChainToken string
}

var tiktok_regexp = regexp.MustCompile(`<script id=\"SIGI_STATE\" type=\"application/json\">(.*?)</script>`)

func NewTikTok(sessionID, chainToken string) TikTok {
	return TikTok{sessionID, chainToken}
}

func (tiktok *TikTok) Post(owner, post string, incognito bool) (URL string, username string, err error) {
	postURL := fmt.Sprintf("https://www.tiktok.com/@%s/video/%s", owner, post)
	request, err := http.NewRequest(http.MethodGet, postURL, nil)
	if err != nil {
		return URL, username, err
	}

	if !incognito {
		sessionCookie := http.Cookie{
			Name:     "sessionid",
			Value:    tiktok.SessionID,
			Domain:   ".tiktok.com",
			HttpOnly: true,
		}
		request.AddCookie(&sessionCookie)
		chainCookie := http.Cookie{
			Name:     "tt_chain_token",
			Value:    tiktok.ChainToken,
			Domain:   ".tiktok.com",
			HttpOnly: true,
			Secure:   true,
		}
		request.AddCookie(&chainCookie)
	}
	request.Header.Add("User-Agent", UserAgent)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return URL, username, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return URL, username, err
	}

	script := tiktok_regexp.FindString(string(body))
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
