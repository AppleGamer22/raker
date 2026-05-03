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

type InstagramHighlight struct {
	Edges [1]struct {
		Node struct {
			User struct {
				Username string `json:"username"`
			} `json:"user"`
			Items []InstagramItem `json:"items"`
		} `json:"node"`
	} `json:"edges"`
}

func (instagram *Instagram) Highlights(id string) ([]string, string, error) {
	htmlRequest, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.instagram.com/stories/highlights/%s", id), nil)
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

	var instagramStory InstagramHighlight
	if err := json.Unmarshal([]byte(reelsMediaJSON), &instagramStory); err != nil {
		return []string{}, "", err
	}

	if len(instagramStory.Edges[0].Node.Items) == 0 {
		return []string{}, "", errors.New("could not find reels_media entries")
	}

	username := instagramStory.Edges[0].Node.User.Username
	URLs := make([]string, 0, len(instagramStory.Edges[0].Node.Items))
	for _, item := range instagramStory.Edges[0].Node.Items {
		URLs = append(URLs, item.URLs()...)
	}

	return URLs, username, nil
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
var storyReelsMediaRegex = regexp.MustCompile(`"xdt_api__v1__feed__reels_media.*?"\s*:\s*`)

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
	URLs := make([]string, 0, len(instagramStory.ReelsMedia[0].Items))
	for _, item := range instagramStory.ReelsMedia[0].Items {
		URLs = append(URLs, item.URLs()...)
	}

	return URLs, username, nil
}
