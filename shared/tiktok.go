package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"slices"

	"github.com/charmbracelet/log"
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
						Cover          string `json:"cover"`
						PlayAddress    string `json:"playAddr"`
						PlayAddrStruct struct {
							UrlList []string
						}
					} `json:"video"`
				} `json:"itemStruct"`
			} `json:"itemInfo"`
		} `json:"webapp.video-detail"`
	} `json:"__DEFAULT_SCOPE__"`
}

type TikTok struct {
	SessionID      string
	SessionIDGuard string
	// ChainToken     string
}

var tiktok_regexp = regexp.MustCompile(`<script id=\"__UNIVERSAL_DATA_FOR_REHYDRATION__\" type=\"application/json\">(.*?)</script>`)

type WAFCookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewTikTok(sessionID, sessionIDGuard string) TikTok {
	return TikTok{sessionID, sessionIDGuard}
}

func (tikok *TikTok) MSToken(owner string) (*http.Client, error) {
	client := NewClient(false)
	ownerURL := fmt.Sprintf("https://www.tiktok.com/@%s", owner)
	request, err := http.NewRequest(http.MethodGet, ownerURL, nil)
	if err != nil {
		return &http.Client{}, err
	}

	response, err := client.Do(request)
	if err != nil {
		return &http.Client{}, err
	}
	response.Body.Close()

	request, err = http.NewRequest(http.MethodPost, "https://mssdk-sg.tiktok.com/web/common", nil)
	if err != nil {
		return &http.Client{}, err
	}

	query := request.URL.Query()
	for _, cookie := range response.Cookies() {
		if cookie.Name == "msToken" {
			query.Add("msToken", cookie.Value)
		}
	}
	request.URL.RawQuery = query.Encode()
	for _, cookie := range request.Cookies() {
		request.AddCookie(cookie)
	}

	response, err = client.Do(request)
	if err != nil {
		return &http.Client{}, err
	}
	response.Body.Close()

	return client, nil
}

func (tiktok *TikTok) FetchPost(owner, post string, incognito bool, wafCookies ...*http.Cookie) (string, []*http.Cookie, error) {
	postURL := fmt.Sprintf("https://www.tiktok.com/@%s/video/%s", owner, post)
	request, err := http.NewRequest(http.MethodGet, postURL, nil)
	if err != nil {
		return "", []*http.Cookie{}, err
	}

	for _, cookie := range wafCookies {
		request.AddCookie(cookie)
	}

	if !incognito {
		sessionCookie := http.Cookie{
			Name:     "sessionid",
			Value:    tiktok.SessionID,
			Domain:   ".tiktok.com",
			HttpOnly: true,
		}
		request.AddCookie(&sessionCookie)
		// chainCookie := http.Cookie{
		// 	Name:     "tt_chain_token",
		// 	Value:    tiktok.ChainToken,
		// 	Domain:   ".tiktok.com",
		// 	HttpOnly: true,
		// 	Secure:   true,
		// }
		// request.AddCookie(&chainCookie)
		sessionGuardCookie := http.Cookie{
			Name:     "sid_guard",
			Value:    tiktok.SessionIDGuard,
			Domain:   ".tiktok.com",
			HttpOnly: true,
		}
		request.AddCookie(&sessionGuardCookie)
	}
	client, err := tiktok.MSToken(owner)
	if err != nil {
		return "", []*http.Cookie{}, err
	}
	// request.Header.Add("sec-ch-ua", `"Not_A Brand";v="99", "Google Chrome";v="109", "Chromium";v="109"`)

	response, err := client.Do(request)
	if err != nil {
		return "", []*http.Cookie{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", []*http.Cookie{}, err
	}

	script := tiktok_regexp.FindString(string(body))
	if script == "" {
		// fmt.Println(string(body))
		return string(body), []*http.Cookie{}, errors.New("could not find JSON")
	}
	return script, response.Cookies(), nil
}

func (tiktok *TikTok) TryWAF(url, challenge string) ([]*http.Cookie, error) {
	command := exec.Command("deno", "run", "--allow-net", "waf.ts", url, challenge)

	output, err := command.Output()
	if err != nil {
		log.Error(err)
		// log.Error(output)
		return []*http.Cookie{}, err
	}

	var wafCookies []WAFCookie
	if err := json.Unmarshal(output, &wafCookies); err != nil {
		log.Error(err)
		return []*http.Cookie{}, err
	} else if len(wafCookies) == 0 {
		return []*http.Cookie{}, errors.New("No WAF cookies were found")
	}

	cookies := make([]*http.Cookie, 0, len(wafCookies))

	for _, wafCookie := range wafCookies {
		cookies = append(cookies, &http.Cookie{
			Name:   wafCookie.Name,
			Value:  wafCookie.Value,
			MaxAge: 10,
		})
	}

	return cookies, nil
}

func (tiktok *TikTok) Post(owner, post string, incognito bool) ([]string, []string, string, []*http.Cookie, error) {

	postURL := fmt.Sprintf("https://www.tiktok.com/@%s/video/%s", owner, post)
	request, err := http.NewRequest(http.MethodGet, postURL, nil)
	if err != nil {
		return []string{}, []string{}, "", []*http.Cookie{}, err
	}

	script, cookies, err := tiktok.FetchPost(owner, post, incognito)
	if err != nil {
		log.Warn("attempting WAF", "owner", owner, "post", post)
		wafCookies, err := tiktok.TryWAF(postURL, script)
		if err != nil {
			log.Error(err)
			return []string{}, []string{}, "", []*http.Cookie{}, errors.New("could not pass WAF")
		}
		script, cookies, err = tiktok.FetchPost(owner, post, incognito, wafCookies...)
		if err != nil {
			log.Error(err)
			return []string{}, []string{}, "", []*http.Cookie{}, errors.New("could not find JSON after WAF")
		}
	}

	jsonText := script[len(`<script id="__UNIVERSAL_DATA_FOR_REHYDRATION__" type="application/json">`) : len(script)-len("</script>")]
	var tiktokPost TikTokPost
	if err := json.Unmarshal([]byte(jsonText), &tiktokPost); err != nil {
		return []string{}, []string{}, "", []*http.Cookie{}, err
	}

	// fmt.Println(jsonText)

	username := tiktokPost.DefaultScop.VideoDetail.ItemInfo.ItemStruct.Author.UniqueID
	URL := tiktokPost.DefaultScop.VideoDetail.ItemInfo.ItemStruct.Video.PlayAddress
	if URL == "" && len(tiktokPost.DefaultScop.VideoDetail.ItemInfo.ItemStruct.Video.PlayAddrStruct.UrlList) == 0 {
		if len(tiktokPost.DefaultScop.VideoDetail.ItemInfo.ItemStruct.ImagePost.Images) == 0 {
			return []string{}, []string{}, "", []*http.Cookie{}, errors.New("Post not available from incognito mode")
		}
		URLs := make([]string, 0, len(tiktokPost.DefaultScop.VideoDetail.ItemInfo.ItemStruct.ImagePost.Images))
		for _, image := range tiktokPost.DefaultScop.VideoDetail.ItemInfo.ItemStruct.ImagePost.Images {
			URLs = append(URLs, image.ImageURL.URLs[0])
		}
		return []string{}, URLs, username, slices.Concat(request.Cookies(), cookies), nil
	} else if URL != "" && len(tiktokPost.DefaultScop.VideoDetail.ItemInfo.ItemStruct.Video.PlayAddrStruct.UrlList) == 0 {
		return []string{URL}, []string{tiktokPost.DefaultScop.VideoDetail.ItemInfo.ItemStruct.Video.Cover}, username, slices.Concat(request.Cookies(), cookies), nil
	}

	return tiktokPost.DefaultScop.VideoDetail.ItemInfo.ItemStruct.Video.PlayAddrStruct.UrlList, []string{tiktokPost.DefaultScop.VideoDetail.ItemInfo.ItemStruct.Video.Cover}, username, slices.Concat(request.Cookies(), cookies), nil
}
