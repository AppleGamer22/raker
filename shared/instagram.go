package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type InstagramUserID struct {
	Data struct {
		User struct {
			ID string `json:"id"`
		} `json:"user"`
	} `json:"data"`
}

type InstagramItem struct {
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
}

func (item *InstagramItem) URLs() []string {
	output := []string{}
	if len(item.CarouselMedia) > 0 {
		for _, media := range item.CarouselMedia {
			if len(media.VideoVersions) > 0 {
				output = append(output, media.VideoVersions[0].URL)
			}
			if len(media.ImageVersions2.Candidates) > 0 {
				output = append(output, media.ImageVersions2.Candidates[0].URL)
			}
		}
	} else {
		if len(item.VideoVersions) > 0 {
			output = append(output, item.VideoVersions[0].URL)
		}
		if len(item.ImageVersions2.Candidates) > 0 {
			output = append(output, item.ImageVersions2.Candidates[0].URL)
		}
	}
	return output
}

type InstagramPost struct {
	Items [1]InstagramItem `json:"items"`
}

type InstagramReels struct {
	ReelsMedia [1]struct {
		User struct {
			Username string `json:"username"`
		} `json:"user"`
		Items []InstagramItem `json:"items"`
	} `json:"reels_media"`
}

type Instagram struct {
	fbsrCookie    http.Cookie
	sessionCookie http.Cookie
	userCookie    http.Cookie
}

var instagram_regexp = regexp.MustCompile(`\"media_id\":\"?([0-9]+)\"?`)

func NewInstagram(fbsr, sessionID, userID string) Instagram {
	return Instagram{
		fbsrCookie: http.Cookie{
			Name:     "fbsr_124024574287414",
			Value:    fbsr,
			Domain:   ".instagram.com",
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
		},
		sessionCookie: http.Cookie{
			Name:     "sessionid",
			Value:    sessionID,
			Domain:   ".instagram.com",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		},
		userCookie: http.Cookie{
			Name:   "ds_user_id",
			Value:  userID,
			Domain: ".instagram.com",
			Path:   "/",
			Secure: true,
		},
	}
}

func (instagram *Instagram) Post(post string, incognito bool) (URLs []string, username string, err error) {
	htmlURL := fmt.Sprintf("https://www.instagram.com/p/%s", post)
	htmlRequest, err := http.NewRequest(http.MethodGet, htmlURL, nil)
	if err != nil {
		return URLs, username, err
	}

	if !incognito {
		htmlRequest.AddCookie(&instagram.fbsrCookie)
		htmlRequest.AddCookie(&instagram.sessionCookie)
		htmlRequest.AddCookie(&instagram.userCookie)
	}

	htmlRequest.Header.Add("x-ig-app-id", "936619743392459")
	htmlRequest.Header.Add("user-agent", UserAgent)
	htmlRequest.Header.Add("referer", "https://www.instagram.com/")
	htmlRequest.Header.Add("sec-fetch-mode", "navigate")

	htmlResponse, err := http.DefaultClient.Do(htmlRequest)
	if err != nil {
		return URLs, username, err
	}
	defer htmlResponse.Body.Close()

	htmlBody, err := io.ReadAll(htmlResponse.Body)
	if err != nil {
		return URLs, username, err
	}

	mediaIDMatch := instagram_regexp.FindString(string(htmlBody))
	if mediaIDMatch == "" {
		return URLs, username, errors.New("could not find media ID")
	}

	mediaID := func() string {
		if mediaIDMatch[len(mediaIDMatch)-1] != '"' {
			return mediaIDMatch[len(`"media_id":`):]
		} else {
			return mediaIDMatch[len(`"media_id":"`) : len(mediaIDMatch)-1]
		}
	}()

	jsonURL := fmt.Sprintf("https://i.instagram.com/api/v1/media/%s/info/", mediaID)
	jsonRequest, err := http.NewRequest(http.MethodGet, jsonURL, nil)
	if err != nil {
		return URLs, username, err
	}

	if !incognito {
		jsonRequest.AddCookie(&instagram.fbsrCookie)
		jsonRequest.AddCookie(&instagram.sessionCookie)
		jsonRequest.AddCookie(&instagram.userCookie)
	}

	jsonRequest.Header.Add("x-ig-app-id", "936619743392459")
	jsonRequest.Header.Add("User-Agent", UserAgent)
	jsonRequest.Header.Add("referer", "https://www.instagram.com/")

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
