package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
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

	// request.AddCookie(&instagram.fbsrCookie)
	request.AddCookie(&instagram.sessionCookie)
	request.AddCookie(&instagram.userCookie)
	// request.Header.Add("x-ig-app-id", "936619743392459")
	request.Header.Add("User-Agent", UserAgent)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return URLs, username, err
	}
	defer response.Body.Close()

	statusClass := response.StatusCode / 100
	if statusClass == 4 || statusClass == 5 {
		return []string{}, "", fmt.Errorf("response of %d instead of media", response.StatusCode)
	}

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

	// request.AddCookie(&instagram.fbsrCookie)
	request.AddCookie(&instagram.sessionCookie)
	request.AddCookie(&instagram.userCookie)
	request.Header.Add("x-ig-app-id", "936619743392459")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var instagramUserID InstagramUserID
	err = json.NewDecoder(response.Body).Decode(&instagramUserID)
	return instagramUserID.Data.User.ID, err
}

type InstagramStory struct {
	ReelsMedia []struct {
		User struct {
			Username string `json:"username"`
		} `json:"user"`
		Items []InstagramItem `json:"items"`
	} `json:"reels_media"`
}

var storyReelRegex = regexp.MustCompile(`<script type="application/json" .*? data-sjs>(.*?xdt_api__v1__feed__reels_media.*?)</script>`)
var storyReelsMediaRegex = regexp.MustCompile(`"xdt_api__v1__feed__reels_media"\s*:\s*`)

func extractReelsMediaJSON(jsonText string) (string, error) {
	match := storyReelsMediaRegex.FindStringIndex(jsonText)
	if match == nil {
		return "", errors.New("could not find reels_media JSON")
	}

	start := match[1]
	if start >= len(jsonText) || jsonText[start] != '{' {
		return "", errors.New("could not find reels_media object")
	}

	depth := 0
	for index := start; index < len(jsonText); index++ {
		switch jsonText[index] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return jsonText[start : index+1], nil
			}
		}
	}

	return "", errors.New("unterminated reels_media object")
}

func (instagram *Instagram) Story(username string) ([]string, string, error) {
	var URLs []string

	htmlRequest, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.instagram.com/stories/%s", username), nil)
	if err != nil {
		return []string{}, "", err
	}

	htmlRequest.AddCookie(&instagram.sessionCookie)
	htmlRequest.AddCookie(&instagram.userCookie)

	client := NewClient(false)

	htmlResponse, err := client.Do(htmlRequest)
	if err != nil {
		return []string{}, "", err
	}
	defer htmlResponse.Body.Close()

	body, err := io.ReadAll(htmlResponse.Body)
	if err != nil {
		return []string{}, "", err
	}

	script := storyReelRegex.FindString(string(body))
	if script == "" {
		return []string{}, "", errors.New("could not find JSON")
	}

	jsonText := script[strings.Index(script, "{") : len(script)-len("</script>")]
	reelsMediaJSON, err := extractReelsMediaJSON(jsonText)
	if err != nil {
		return []string{}, "", err
	}

	var instagramStory InstagramStory
	if err := json.Unmarshal([]byte(reelsMediaJSON), &instagramStory); err != nil {
		return []string{}, "", err
	}

	if len(instagramStory.ReelsMedia) == 0 {
		return []string{}, "", errors.New("could not find reels_media entries")
	}

	username = instagramStory.ReelsMedia[0].User.Username
	for _, item := range instagramStory.ReelsMedia[0].Items {
		URLs = append(URLs, item.URLs()...)
	}

	return URLs, username, nil
}
