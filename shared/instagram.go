package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

func (instagram *Instagram) Post(post string) (URLs []string, username string, err error) {
	htmlURL := fmt.Sprintf("https://www.instagram.com/p/%s", post)
	htmlRequest, err := http.NewRequest(http.MethodGet, htmlURL, nil)
	if err != nil {
		return URLs, username, err
	}
	htmlRequest.AddCookie(&instagram.fbsr)
	htmlRequest.AddCookie(&instagram.sessionID)
	htmlRequest.Header.Add("x-ig-app-id", instagram.appID)
	htmlRequest.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")

	htmlResponse, err := http.DefaultClient.Do(htmlRequest)
	if err != nil {
		return URLs, username, err
	}
	defer htmlResponse.Body.Close()

	htmlBody, err := io.ReadAll(htmlResponse.Body)
	if err != nil {
		return URLs, username, err
	}
	re := regexp.MustCompile(`media\?id=([0-9]+)`)
	mediaIDMatch := re.FindString(string(htmlBody))
	if mediaIDMatch == "" {
		fmt.Println(string(htmlBody))
		return URLs, username, errors.New("could not find media ID")
	}

	mediaID := mediaIDMatch[len(`media?id=`):]
	jsonURL := fmt.Sprintf("https://i.instagram.com/api/v1/media/%s/info/", mediaID)
	jsonRequest, err := http.NewRequest(http.MethodGet, jsonURL, nil)
	if err != nil {
		return URLs, username, err
	}

	jsonRequest.AddCookie(&instagram.fbsr)
	jsonRequest.AddCookie(&instagram.sessionID)
	jsonRequest.Header.Add("x-ig-app-id", instagram.appID)
	jsonRequest.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")

	jsonResponse, err := http.DefaultClient.Do(jsonRequest)
	if err != nil {
		return URLs, username, err
	}
	defer jsonResponse.Body.Close()

	var instagramPost InstagramPost
	if err := json.NewDecoder(jsonResponse.Body).Decode(&instagramPost); err != nil {
		return URLs, username, err
	}

	item := instagramPost.Items[0]
	username = item.User.Username
	URLs = item.URLs()

	return URLs, username, err
}
