package shared

import (
	"fmt"

	"github.com/chromedp/chromedp"
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
	EntryData struct {
		PostPage []struct {
			GraphQL struct {
				ShortcodeMedia struct {
					DisplayURL            string `json:"display_url"`
					VideoURL              string `json:"video_url"`
					EdgeSidecarToChildren struct {
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
				} `json:"shortcode_media"`
			} `json:"graphql"`
		}
	} `json:"entry_data"`
}

func (browser Browser) InstagramPost(post string) (URLs []string, username string, err error) {
	defer browser.CannelAllocator()
	defer browser.CancelTask()
	postURL := fmt.Sprintf("https://www.instagram.com/p/%s", post)

	var instagramPost InstagramPost

	err = chromedp.Run(browser.Task,
		chromedp.Navigate(postURL),
		chromedp.WaitReady(browser.InstagramScriptSelector()),
		chromedp.Evaluate(browser.InstagramScript(post), &instagramPost),
	)

	if err != nil {
		return URLs, username, err
	}

	if browser.Incognito {
		page := instagramPost.EntryData.PostPage[0]
		username = page.GraphQL.ShortcodeMedia.Owner.Username
		if len(page.GraphQL.ShortcodeMedia.EdgeSidecarToChildren.Edges) > 0 {
			for _, edge := range page.GraphQL.ShortcodeMedia.EdgeSidecarToChildren.Edges {
				if edge.Node.VideoURL != "" {
					URLs = append(URLs, edge.Node.VideoURL)
				} else {
					URLs = append(URLs, edge.Node.DisplayURL)
				}
			}
		} else {
			if page.GraphQL.ShortcodeMedia.VideoURL != "" {
				URLs = append(URLs, page.GraphQL.ShortcodeMedia.VideoURL)
			} else {
				URLs = append(URLs, page.GraphQL.ShortcodeMedia.DisplayURL)
			}
		}
	} else {
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
	}

	return URLs, username, err
}

func (browser Browser) InstagramScriptSelector() string {
	if browser.Debug || browser.Incognito {
		return "script:nth-child(15)"
	}
	return "script:nth-child(16)"
}

func (browser Browser) InstagramScript(post string) string {
	prefixLength := len("window.__additionalDataLoaded(/p/") + len(post) + 4
	if browser.Incognito {
		return "window._sharedData"
	} else if browser.Debug {
		return fmt.Sprintf(`JSON.parse(document.querySelector("script:nth-child(15)").text.slice(%d, -2))`, prefixLength)
	}
	return fmt.Sprintf(`JSON.parse(document.querySelector("script:nth-child(16)").text.slice(%d, -2))`, prefixLength)
}
