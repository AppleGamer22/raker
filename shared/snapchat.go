package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type SnapchatHighlight struct {
	Props struct {
		PageProps struct {
			PublicUserProfile struct {
				Username string `json:"username"`
			} `json:"publicUserProfile"`
			Highlight struct {
				SnapList []struct {
					// 0=image | 1=video
					SnapMediaType int `json:"snapMediaType"`
					SnapURLs      struct {
						MediaURL string `json:"mediaUrl"`
					} `json:"snapUrls"`
				} `json:"snapList"`
			} `json:"highlight"`
		} `json:"pageProps"`
	} `json:"props"`
}

type SnapchatHighlightResult struct {
	Username string
	URLs     []struct {
		IsVideo bool
		URL     string
	}
}

var snapchatRegex = regexp.MustCompile(`<script id="__NEXT_DATA__" type="application/json">(.*?)</script>`)

func Snapchat(owner, highlight string) (SnapchatHighlightResult, []*http.Cookie, error) {
	postURL := fmt.Sprintf("https://www.snapchat.com/@%s/highlight/%s", owner, highlight)

	client := NewClient(false)

	htmlRequest, err := http.NewRequest(http.MethodGet, postURL, nil)
	if err != nil {
		return SnapchatHighlightResult{}, []*http.Cookie{}, err
	}

	htmlResponse, err := client.Do(htmlRequest)
	if err != nil {
		return SnapchatHighlightResult{}, []*http.Cookie{}, err
	}
	defer htmlResponse.Body.Close()

	body, err := io.ReadAll(htmlResponse.Body)
	if err != nil {
		return SnapchatHighlightResult{}, []*http.Cookie{}, err
	}

	script := snapchatRegex.FindString(string(body))
	if script == "" {
		return SnapchatHighlightResult{}, []*http.Cookie{}, errors.New("could not find JSON")
	}

	jsonText := script[len(`<script id="__NEXT_DATA__" type="application/json">`) : len(script)-len("</script>")]
	var snapchatHighlight SnapchatHighlight
	if err := json.Unmarshal([]byte(jsonText), &snapchatHighlight); err != nil {
		return SnapchatHighlightResult{}, []*http.Cookie{}, err
	}

	var result SnapchatHighlightResult
	result.Username = snapchatHighlight.Props.PageProps.PublicUserProfile.Username
	for _, snap := range snapchatHighlight.Props.PageProps.Highlight.SnapList {
		result.URLs = append(result.URLs, struct {
			IsVideo bool
			URL     string
		}{
			URL:     snap.SnapURLs.MediaURL,
			IsVideo: snap.SnapMediaType == 1,
		})
	}

	return result, []*http.Cookie{}, nil
}
