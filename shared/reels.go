package shared

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (instagram *Instagram) Reels(id string, highlight bool) (URLs []string, username string, err error) {
	request, err := http.NewRequest(http.MethodGet, "https://i.instagram.com/api/v1/feed/reels_media/", nil)
	if err != nil {
		return URLs, username, err
	}

	query := request.URL.Query()
	if highlight {
		query.Add("reel_ids", fmt.Sprintf("highlight:%s", id))
	} else {
		id, err = instagram.userID(id)
		if err != nil {
			return URLs, username, err
		}
		query.Add("reel_ids", id)
	}
	request.URL.RawQuery = query.Encode()

	request.AddCookie(&instagram.fbsr)
	request.AddCookie(&instagram.sessionID)
	request.Header.Add("x-ig-app-id", instagram.appID)
	request.Header.Add("user-agent", userAgent)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return URLs, username, err
	}
	defer response.Body.Close()

	var instagramReels InstagramReels
	if err := json.NewDecoder(response.Body).Decode(&instagramReels); err != nil {
		return URLs, username, err
	}

	username = instagramReels.ReelsMedia[0].User.Username
	for _, item := range instagramReels.ReelsMedia[0].Items {
		URLs = append(URLs, item.URLs()...)
	}

	return URLs, username, err
}

func (instagram *Instagram) userID(username string) (string, error) {
	request, err := http.NewRequest(http.MethodGet, "https://i.instagram.com/api/v1/users/web_profile_info/", nil)
	if err != nil {
		return "", err
	}

	query := request.URL.Query()
	query.Add("username", username)
	request.URL.RawQuery = query.Encode()

	request.AddCookie(&instagram.fbsr)
	request.AddCookie(&instagram.sessionID)
	request.Header.Add("x-ig-app-id", instagram.appID)
	request.Header.Add("user-agent", userAgent)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var instagramUserID InstagramUserID
	err = json.NewDecoder(response.Body).Decode(&instagramUserID)
	return instagramUserID.Data.User.ID, err
}
