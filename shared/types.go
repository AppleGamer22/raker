package shared

import "net/http"

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
			} else if len(media.ImageVersions2.Candidates) > 0 {
				output = append(output, media.ImageVersions2.Candidates[0].URL)
			}
		}
	} else {
		if len(item.VideoVersions) > 0 {
			output = append(output, item.VideoVersions[0].URL)
		} else if len(item.ImageVersions2.Candidates) > 0 {
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
	fbsr      http.Cookie
	sessionID http.Cookie
	appID     string
}

func NewInstagram(fbsr, sessionID, appID string) Instagram {
	return Instagram{
		fbsr: http.Cookie{
			Name:     "fbsr_124024574287414",
			Value:    fbsr,
			Domain:   ".instagram.com",
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
		},
		sessionID: http.Cookie{
			Name:     "sessionid",
			Value:    sessionID,
			Domain:   ".instagram.com",
			Path:     "/",
			HttpOnly: true,
		},
		appID: appID,
	}
}
