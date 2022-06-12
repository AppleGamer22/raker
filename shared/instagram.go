package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type InstagramPost struct {
	Items []struct {
		CarouselMedia []struct {
			ImageVersions2 struct {
				Candidates []struct {
					URL string `json:"url"`
				} `json:"candidates"`
			} `json:"image_versions2"`
			VideoVersions []struct {
				URL string `json:"url"`
			} `json:"video_versions"`
		} `json:"carousel_media"`
		ImageVersions2 struct {
			Candidates []struct {
				URL string `json:"url"`
			} `json:"candidates"`
		} `json:"image_versions2"`
		VideoVersions []struct {
			URL string `json:"url"`
		} `json:"video_versions"`
		User struct {
			Username string `json:"username"`
		} `json:"user"`
	} `json:"items"`
}

type Instagram struct {
	fbsr      string
	sessionID string
	appID     string
}

func NewInstagram(fbsr, sessionID, appID string) Instagram {
	return Instagram{
		fbsr:      fbsr,
		sessionID: sessionID,
		appID:     appID,
	}
}

func (instagram *Instagram) Do(post string) (URLs []string, username string, err error) {
	fbsrCookie := http.Cookie{
		Name:     "fbsr_124024574287414",
		Value:    instagram.fbsr,
		Domain:   ".instagram.com",
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	}

	sessionCookie := http.Cookie{
		Name:     "sessionid",
		Value:    instagram.sessionID,
		Domain:   ".instagram.com",
		Path:     "/",
		HttpOnly: true,
	}

	htmlURL := fmt.Sprintf("https://www.instagram.com/p/%s", post)
	htmlRequest, err := http.NewRequest(http.MethodGet, htmlURL, nil)
	if err != nil {
		return URLs, username, err
	}
	htmlRequest.AddCookie(&fbsrCookie)
	htmlRequest.AddCookie(&sessionCookie)
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

	jsonRequest.AddCookie(&fbsrCookie)
	jsonRequest.AddCookie(&sessionCookie)
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
	if len(item.CarouselMedia) > 0 {
		for _, media := range item.CarouselMedia {
			if len(media.VideoVersions) > 0 {
				URLs = append(URLs, media.VideoVersions[0].URL)
			} else if len(media.ImageVersions2.Candidates) > 0 {
				URLs = append(URLs, media.ImageVersions2.Candidates[0].URL)
			}
		}
	} else {
		if len(item.VideoVersions) > 0 {
			URLs = append(URLs, item.VideoVersions[0].URL)
		} else if len(item.ImageVersions2.Candidates) > 0 {
			URLs = append(URLs, item.ImageVersions2.Candidates[0].URL)
		}
	}

	return URLs, username, err
}
