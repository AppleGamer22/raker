package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type InstagramUserID struct {
	Data struct {
		User struct {
			ID string `json:"id"`
		} `json:"user"`
	} `json:"data"`
}

type InstagramPostIncognito struct {
	Data struct {
		ShortcodeMedia struct {
			EdgeSidecarChildren struct {
				Edges []struct {
					Node struct {
						DisplayURL string `json:"display_url"`
						VideoURL   string `json:"video_url"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_sidecar_to_children"`
			Owner struct {
				Username string `json:"username"`
			} `json:"owner"`
			DisplayURL string `json:"display_url"`
			VideoURL   string `json:"video_url"`
		} `json:"xdt_shortcode_media"`
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
	output := make([]string, 0, len(item.CarouselMedia)+1)
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

func (post *InstagramPostIncognito) URLs() []string {
	output := make([]string, 0, len(post.Data.ShortcodeMedia.EdgeSidecarChildren.Edges)+1)
	if len(post.Data.ShortcodeMedia.EdgeSidecarChildren.Edges) > 0 {
		for _, media := range post.Data.ShortcodeMedia.EdgeSidecarChildren.Edges {
			if media.Node.VideoURL != "" {
				output = append(output, media.Node.VideoURL)
			}
			if media.Node.DisplayURL != "" {
				output = append(output, media.Node.DisplayURL)
			}
		}
	} else {
		if post.Data.ShortcodeMedia.VideoURL != "" {
			output = append(output, post.Data.ShortcodeMedia.VideoURL)
		}
		if post.Data.ShortcodeMedia.DisplayURL != "" {
			output = append(output, post.Data.ShortcodeMedia.DisplayURL)
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

var (
	instagramRegExpMediaID              = regexp.MustCompile(`media_id\":\"([0-9]+)`)
	instagramRegExpDATR                 = regexp.MustCompile(`_js_datr\":{\"value":\"([0-9a-zA-Z-]+)`)
	instagramRegExpLSD                  = regexp.MustCompile(`lsd\":\"([0-9a-zA-Z-]+)`)
	instagramRegExpScriptWithDocumentID = regexp.MustCompile(`<link rel=\"preload\" href=\"(.*?)\" as=\"script\" crossorigin=\"anonymous\" nonce=".*?" />`)
	instagramRegExpDocumentID           = regexp.MustCompile(`params:{id:\"([0-9]+)\",metadata:{},name:\"PolarisPostActionLoadPostQueryQuery`)
)

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

func (instagram *Instagram) Post(post string) (URLs []string, username string, err error) {
	htmlURL := fmt.Sprintf("https://www.instagram.com/p/%s", post)
	htmlRequest, err := http.NewRequest(http.MethodGet, htmlURL, nil)
	if err != nil {
		return URLs, username, err
	}

	htmlRequest.AddCookie(&instagram.fbsrCookie)
	htmlRequest.AddCookie(&instagram.sessionCookie)
	htmlRequest.AddCookie(&instagram.userCookie)

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

	mediaIDMatch := instagramRegExpMediaID.FindString(string(htmlBody))
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

	jsonRequest.AddCookie(&instagram.fbsrCookie)
	jsonRequest.AddCookie(&instagram.sessionCookie)
	jsonRequest.AddCookie(&instagram.userCookie)

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

func InstagramIncognito(post string) ([]string, string, []*http.Cookie, error) {
	htmlURL := fmt.Sprintf("https://www.instagram.com/p/%s", post)
	htmlRequest, err := http.NewRequest(http.MethodGet, htmlURL, nil)
	if err != nil {
		return []string{}, "", []*http.Cookie{}, err
	}

	htmlRequest.Header.Add("x-ig-app-id", "936619743392459")
	htmlRequest.Header.Add("user-agent", UserAgent)
	htmlRequest.Header.Add("referer", "https://www.instagram.com/")
	htmlRequest.Header.Add("sec-fetch-mode", "navigate")

	htmlResponse, err := http.DefaultClient.Do(htmlRequest)
	if err != nil {
		return []string{}, "", []*http.Cookie{}, err
	}
	defer htmlResponse.Body.Close()

	htmlBody, err := io.ReadAll(htmlResponse.Body)
	if err != nil {
		return []string{}, "", []*http.Cookie{}, err
	}

	datrMatches := instagramRegExpDATR.FindStringSubmatch(string(htmlBody))
	if len(datrMatches) == 0 || datrMatches[1] == "" {
		return []string{}, "", []*http.Cookie{}, errors.New("could not find datr value")
	}

	lsdMatches := instagramRegExpLSD.FindStringSubmatch(string(htmlBody))
	if len(datrMatches) == 0 || datrMatches[1] == "" {
		return []string{}, "", []*http.Cookie{}, errors.New("could not find lsd value")
	}

	jsURLs := instagramRegExpScriptWithDocumentID.FindAllStringSubmatch(string(htmlBody), 4)
	if jsURLs == nil || jsURLs[3][1] == "" {
		return []string{}, "", []*http.Cookie{}, errors.New("could not find link URL")
	}
	jsURL := jsURLs[3][1]

	jsRequest, err := http.NewRequest(http.MethodGet, jsURL, nil)
	if err != nil {
		return []string{}, "", []*http.Cookie{}, err
	}
	jsRequest.Header.Add("user-agent", UserAgent)
	jsRequest.Header.Add("referer", "https://www.instagram.com/")

	jsResponse, err := http.DefaultClient.Do(jsRequest)
	if err != nil {
		return []string{}, "", []*http.Cookie{}, err
	}
	defer jsResponse.Body.Close()

	jsBody, err := io.ReadAll(jsResponse.Body)
	if err != nil {
		return []string{}, "", []*http.Cookie{}, err
	}
	documentIDs := instagramRegExpDocumentID.FindStringSubmatch(string(jsBody))
	if documentIDs == nil || documentIDs[1] == "" {
		return []string{}, "", []*http.Cookie{}, errors.New("could not find document ID")
	}

	form := url.Values{
		"lsd":       {lsdMatches[1]},
		"doc_id":    {documentIDs[1]},
		"variables": {fmt.Sprintf(`{"shortcode":"%s","fetch_comment_count":40,"fetch_related_profile_media_count":3,"parent_comment_count":24,"child_comment_count":3,"fetch_like_count":10,"fetch_tagged_user_count":null,"fetch_preview_comment_count":2,"has_threaded_comments":true,"hoisted_comment_id":null,"hoisted_reply_id":null}`, post)},
	}
	jsonRequest, err := http.NewRequest(http.MethodPost, "https://www.instagram.com/api/graphql", strings.NewReader(form.Encode()))
	if err != nil {
		return []string{}, "", []*http.Cookie{}, err
	}
	jsonRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	jsonRequest.Header.Add("sec-fetch-site", "same-origin")
	jsonRequest.Header.Add("x-ig-app-id", "936619743392459")
	jsonRequest.Header.Add("x-asbd-id", "129477")
	jsonRequest.Header.Add("x-fb-friendly-name", "PolarisPostActionLoadPostQueryQuery")
	jsonRequest.Header.Add("User-Agent", UserAgent)
	jsonRequest.Header.Add("referer", htmlURL)

	jsonRequest.AddCookie(&http.Cookie{
		Domain:   ".instagram.com",
		Name:     "datr",
		Value:    datrMatches[1],
		Path:     "/",
		Expires:  time.Now().AddDate(2, 0, 0),
		Secure:   true,
		HttpOnly: true,
	})
	for _, cookie := range htmlResponse.Cookies() {
		jsonRequest.AddCookie(cookie)
	}

	jsonResponse, err := http.DefaultClient.Do(jsonRequest)
	if err != nil {
		return []string{}, "", []*http.Cookie{}, err
	}
	defer jsonResponse.Body.Close()

	// body, _ := io.ReadAll(jsonResponse.Body)
	// fmt.Println(string(body))

	var instagramPost InstagramPostIncognito
	if err := json.NewDecoder(jsonResponse.Body).Decode(&instagramPost); err != nil {
		return []string{}, "", []*http.Cookie{}, err
	}

	username := instagramPost.Data.ShortcodeMedia.Owner.Username
	URLs := instagramPost.URLs()

	return URLs, username, htmlRequest.Cookies(), nil
}
