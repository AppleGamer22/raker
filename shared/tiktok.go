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
	DefaultScop struct {
		VideoDetail struct {
			ItemInfo struct {
				ItemStruct struct {
					Author struct {
						UniqueID string `json:"uniqueId"`
					} `json:"author"`
					ImagePost struct {
						Images []struct {
							ImageURL struct {
								URLs [1]string `json:"urlList"`
							} `json:"imageURL"`
						} `json:"images"`
					} `json:"imagePost"`
					Video struct {
						PlayAddress string `json:"playAddr"`
					} `json:"video"`
				} `json:"itemStruct"`
			} `json:"itemInfo"`
		} `json:"webapp.video-detail"`
	} `json:"__DEFAULT_SCOPE__"`
}

type TikTok struct {
	SessionID      string
	SessionIDGuard string
	ChainToken     string
}

var tiktok_regexp = regexp.MustCompile(`<script id=\"__UNIVERSAL_DATA_FOR_REHYDRATION__\" type=\"application/json\">(.*?)</script>`)

func NewTikTok(sessionID, sessionIDGuard, chainToken string) TikTok {
	return TikTok{sessionID, sessionIDGuard, chainToken}
}

func (tiktok *TikTok) Post(owner, post string, incognito bool) ([]string, string, []*http.Cookie, error) {
	postURL := fmt.Sprintf("https://www.tiktok.com/@%s/video/%s", owner, post)
	request, err := http.NewRequest(http.MethodGet, postURL, nil)
	if err != nil {
		return []string{}, "", []*http.Cookie{}, err
	}

	// if !incognito {
	// 	sessionCookie := http.Cookie{
	// 		Name:     "sessionid",
	// 		Value:    tiktok.SessionID,
	// 		Domain:   ".tiktok.com",
	// 		HttpOnly: true,
	// 	}
	// 	request.AddCookie(&sessionCookie)
	// 	chainCookie := http.Cookie{
	// 		Name:     "tt_chain_token",
	// 		Value:    tiktok.ChainToken,
	// 		Domain:   ".tiktok.com",
	// 		HttpOnly: true,
	// 		Secure:   true,
	// 	}
	// 	request.AddCookie(&chainCookie)
	// 	sessionGuardCookie := http.Cookie{
	// 		Name:     "sid_guard",
	// 		Value:    tiktok.SessionIDGuard,
	// 		Domain:   ".tiktok.com",
	// 		HttpOnly: true,
	// 	}
	// 	request.AddCookie(&sessionGuardCookie)
	// }
	request.Header.Add("User-Agent", UserAgent)
	// request.Header.Add("sec-ch-ua", `"Not_A Brand";v="99", "Google Chrome";v="109", "Chromium";v="109"`)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return []string{}, "", []*http.Cookie{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []string{}, "", []*http.Cookie{}, err
	}

	script := tiktok_regexp.FindString(string(body))
	if script == "" {
		return []string{}, "", []*http.Cookie{}, errors.New("could not find JSON")
	}

	jsonText := script[len(`<script id="__UNIVERSAL_DATA_FOR_REHYDRATION__" type="application/json">`) : len(script)-len("</script>")]
	var tiktokPost TikTokPost
	if err := json.Unmarshal([]byte(jsonText), &tiktokPost); err != nil {
		return []string{}, "", []*http.Cookie{}, err
	}

	username := tiktokPost.DefaultScop.VideoDetail.ItemInfo.ItemStruct.Author.UniqueID
	URL := tiktokPost.DefaultScop.VideoDetail.ItemInfo.ItemStruct.Video.PlayAddress
	if URL == "" {
		URLs := make([]string, 0, len(tiktokPost.DefaultScop.VideoDetail.ItemInfo.ItemStruct.ImagePost.Images))
		for _, image := range tiktokPost.DefaultScop.VideoDetail.ItemInfo.ItemStruct.ImagePost.Images {
			URLs = append(URLs, image.ImageURL.URLs[0])
		}
		return URLs, username, response.Cookies(), err
	}

	return []string{URL}, username, response.Cookies(), err
}
