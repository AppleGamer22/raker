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

type Snap struct {
	// 0=image | 1=video
	SnapMediaType int `json:"snapMediaType"`
	SnapURLs      struct {
		MediaURL string `json:"mediaUrl"`
	} `json:"snapUrls"`
}

type SnapchatHighlightProps struct {
	Props struct {
		PageProps struct {
			PublicUserProfile struct {
				Username string `json:"username"`
			} `json:"publicUserProfile"`
			Highlight struct {
				SnapList []Snap `json:"snapList"`
			} `json:"highlight"`
		} `json:"pageProps"`
	} `json:"props"`
}

type SnapchatStoryProps struct {
	Props struct {
		PageProps struct {
			Story struct {
				SnapList []Snap `json:"snapList"`
				StoryID  struct {
					Value string `json:"value"`
				} `json:"storyId"`
			} `json:"story"`
		} `json:"pageProps"`
	} `json:"props"`
}

type SnapchatHighlightResult struct {
	Username string
	URLs     []struct {
		IsVideo bool
		URL     string
	}
	IsStory bool
}

var (
	snapchatRegex = regexp.MustCompile(`<script id="__NEXT_DATA__" type="application/json">(.*?)</script>`)
	rscPushRegex  = regexp.MustCompile(`__next_f\.push\(\[\s*1\s*,\s*"((?:[^"\\]|\\.)*)"\s*\]\)`)
	userRegex     = regexp.MustCompile(`"username"\s*:\s*"([^"]+)"`)
)

type sequenceItem struct {
	Kind     string `json:"kind"`
	ID       string `json:"id"`
	Src      string `json:"src"`
	MediaURL string `json:"mediaUrl"`
	URL      string `json:"url"`
}

type sequenceData struct {
	Items []sequenceItem `json:"items"`
}

func extractJSONBraces(s string, startIdx int) (string, error) {
	braceStart := strings.Index(s[startIdx:], "{")
	if braceStart == -1 {
		return "", fmt.Errorf("opening brace not found")
	}
	braceStart += startIdx

	depth := 0
	inStr := false
	escape := false
	for i := braceStart; i < len(s); i++ {
		ch := s[i]
		if escape {
			escape = false
			continue
		}
		if ch == '\\' {
			escape = true
			continue
		}
		if ch == '"' {
			inStr = !inStr
			continue
		}
		if !inStr {
			if ch == '{' {
				depth++
			} else if ch == '}' {
				depth--
				if depth == 0 {
					return s[braceStart : i+1], nil
				}
			}
		}
	}
	return "", fmt.Errorf("closing brace not found")
}

func parseRSCMedia(html string, defaultOwner string) (SnapchatHighlightResult, error) {
	matches := rscPushRegex.FindAllStringSubmatch(html, -1)
	if len(matches) == 0 {
		return SnapchatHighlightResult{}, errors.New("no RSC chunks found")
	}

	var stream strings.Builder
	for _, m := range matches {
		var unescaped string
		if err := json.Unmarshal([]byte(`"`+m[1]+`"`), &unescaped); err == nil {
			stream.WriteString(unescaped)
		}
	}
	joined := stream.String()

	seqIdx := strings.Index(joined, `"sequence":`)
	if seqIdx == -1 {
		return SnapchatHighlightResult{}, errors.New("sequence not found in RSC data")
	}

	seqJSON, err := extractJSONBraces(joined, seqIdx)
	if err != nil {
		return SnapchatHighlightResult{}, fmt.Errorf("extracting sequence JSON: %w", err)
	}

	var seq sequenceData
	if err := json.Unmarshal([]byte(seqJSON), &seq); err != nil {
		return SnapchatHighlightResult{}, fmt.Errorf("unmarshaling sequence: %w", err)
	}

	var result SnapchatHighlightResult
	userMatches := userRegex.FindStringSubmatch(joined)
	if len(userMatches) > 1 && userMatches[1] != "" {
		result.Username = userMatches[1]
	} else {
		result.Username = strings.TrimPrefix(defaultOwner, "@")
	}

	for _, item := range seq.Items {
		mediaURL := item.Src
		if mediaURL == "" {
			mediaURL = item.MediaURL
		}
		if mediaURL == "" {
			mediaURL = item.URL
		}
		if mediaURL == "" {
			continue
		}

		result.URLs = append(result.URLs, struct {
			IsVideo bool
			URL     string
		}{
			URL:     mediaURL,
			IsVideo: strings.EqualFold(item.Kind, "video"),
		})
	}

	if len(result.URLs) == 0 {
		return result, errors.New("no media URLs found in sequence")
	}

	return result, nil
}

func SnapchatHighlight(owner, highlight string) (SnapchatHighlightResult, []*http.Cookie, error) {
	postURL := fmt.Sprintf("https://www.snapchat.com/@%s/highlight/%s", owner, highlight)

	client := NewClient([]string{})

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

	bodyStr := string(body)

	// Try RSC stream first (modern Next.js App Router)
	if rscResult, err := parseRSCMedia(bodyStr, owner); err == nil {
		return rscResult, []*http.Cookie{}, nil
	}

	// Fallback to legacy __NEXT_DATA__
	script := snapchatRegex.FindString(bodyStr)
	if script != "" {
		jsonText := script[len(`<script id="__NEXT_DATA__" type="application/json">`) : len(script)-len("</script>")]
		var snapchatHighlight SnapchatHighlightProps
		if err := json.Unmarshal([]byte(jsonText), &snapchatHighlight); err == nil {
			var result SnapchatHighlightResult
			result.Username = snapchatHighlight.Props.PageProps.PublicUserProfile.Username
			if result.Username == "" {
				result.Username = strings.TrimPrefix(owner, "@")
			}
			for _, snap := range snapchatHighlight.Props.PageProps.Highlight.SnapList {
				result.URLs = append(result.URLs, struct {
					IsVideo bool
					URL     string
				}{
					URL:     snap.SnapURLs.MediaURL,
					IsVideo: snap.SnapMediaType == 1,
				})
			}
			if len(result.URLs) > 0 {
				return result, []*http.Cookie{}, nil
			}
		}
	}

	return SnapchatHighlightResult{}, []*http.Cookie{}, errors.New("could not find media data in page")
}

func SnapchatStory(owner string) (SnapchatHighlightResult, []*http.Cookie, error) {
	postURL := fmt.Sprintf("https://www.snapchat.com/@%s", owner)

	client := NewClient([]string{})

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

	bodyStr := string(body)

	// Try RSC stream first (modern Next.js App Router)
	if rscResult, err := parseRSCMedia(bodyStr, owner); err == nil {
		rscResult.IsStory = true
		return rscResult, []*http.Cookie{}, nil
	}

	// Fallback to legacy __NEXT_DATA__
	script := snapchatRegex.FindString(bodyStr)
	if script != "" {
		jsonText := script[len(`<script id="__NEXT_DATA__" type="application/json">`) : len(script)-len("</script>")]
		var snapchatHighlight SnapchatStoryProps
		if err := json.Unmarshal([]byte(jsonText), &snapchatHighlight); err == nil {
			var result SnapchatHighlightResult
			result.Username = snapchatHighlight.Props.PageProps.Story.StoryID.Value
			if result.Username == "" {
				result.Username = strings.TrimPrefix(owner, "@")
			}
			for _, snap := range snapchatHighlight.Props.PageProps.Story.SnapList {
				result.URLs = append(result.URLs, struct {
					IsVideo bool
					URL     string
				}{
					URL:     snap.SnapURLs.MediaURL,
					IsVideo: snap.SnapMediaType == 1,
				})
			}
			result.IsStory = true
			if len(result.URLs) > 0 {
				return result, []*http.Cookie{}, nil
			}
		}
	}

	return SnapchatHighlightResult{}, []*http.Cookie{}, errors.New("could not find media data in page")
}
