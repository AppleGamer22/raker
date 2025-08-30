package shared

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
)

type VSCOPost struct {
	Medias struct {
		ByID map[string]struct {
			Media struct {
				PermaSubdomain string `json:"permaSubdomain"`
				ResponsiveURL  string `json:"responsiveUrl"`
				VideoURL       string `json:"videoUrl"`
				PlaybackURL    string `json:"playbackUrl"`
				Site           struct {
					Domain string `json:"domain"`
				} `json:"site"`
			} `json:"media"`
		} `json:"byId"`
	} `json:"medias"`
}

func findFirstURL(response io.ReadCloser) string {
	scanner := bufio.NewScanner(response)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "https://") {
			return line
		}
	}
	return ""
}

func extractStreamURL(playbackURL string) (string, error) {
	// https://vsco.co/annameadd/media/d939d55c-9543-4b3a-92e5-701d28246e79
	playlistResponse, err := http.Get(playbackURL)
	if err != nil {
		return "", err
	}
	defer playlistResponse.Body.Close()
	renditionURL := findFirstURL(playlistResponse.Body)
	if len(renditionURL) == 0 {
		return "", errors.New("couldn't find rendition URL")
	}
	renditionResponse, err := http.Get(renditionURL)
	if err != nil {
		return "", err
	}
	defer renditionResponse.Body.Close()
	streamURL := findFirstURL(renditionResponse.Body)
	if len(streamURL) == 0 {
		return "", errors.New("couldn't find stream URL")
	}
	return streamURL, nil
}

var vsco_regexp = regexp.MustCompile(`<script>window\.__PRELOADED_STATE__ =(.*?)</script>`)

func VSCO(owner, post string) (string, string, []*http.Cookie, error) {
	postURL := fmt.Sprintf("https://vsco.co/%s/media/%s", owner, post)

	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", "", []*http.Cookie{}, err
	}

	client := &http.Client{
		Jar: jar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS13,
			},
		},
	}

	htmlResponse, err := client.Get(postURL)
	if err != nil {
		return "", "", []*http.Cookie{}, err
	}
	defer htmlResponse.Body.Close()

	body, err := io.ReadAll(htmlResponse.Body)
	if err != nil {
		return "", "", []*http.Cookie{}, err
	}

	script := vsco_regexp.FindString(string(body))
	if script == "" {
		return "", "", []*http.Cookie{}, errors.New("could not find JSON")
	}

	jsonText := script[len("<script>window.__PRELOADED_STATE__ =") : len(script)-len("</script>")]
	jsonText = strings.ReplaceAll(jsonText, "undefined", "null")
	var vscoPost VSCOPost
	if err := json.Unmarshal([]byte(jsonText), &vscoPost); err != nil {
		return "", "", []*http.Cookie{}, err
	}

	media := vscoPost.Medias.ByID[post]
	username := media.Media.PermaSubdomain
	var URL string

	if len(media.Media.VideoURL) > 0 {
		URL = fmt.Sprintf("https://%s", media.Media.VideoURL)
	} else if len(media.Media.PlaybackURL) > 0 {
		username = media.Media.Site.Domain
		URL, err = extractStreamURL(media.Media.PlaybackURL)
	} else {
		URL = fmt.Sprintf("https://%s", media.Media.ResponsiveURL)
	}

	return URL, username, jar.Cookies(htmlResponse.Request.URL), err
}
