package shared

import (
	"fmt"

	"github.com/chromedp/chromedp"
)

const VSCOScriptSelector = "body > script:nth-child(3)"

var VSCOScript = fmt.Sprintf(`JSON.parse(document.querySelector("body > script:nth-child(3)").text.slice(%d))`, len("window.__PRELOADED_STATE__ = "))

type VSCOPost struct {
	Medias struct {
		ByID map[string]struct {
			PermaSubdomain string `json:"permaSubdomain"`
			ResponsiveURL  string `json:"responsiveUrl"`
			VideoURL       string `json:"videoUrl"`
		} `json:"byId"`
	} `json:"medias"`
}

func (browser Browser) VSCO(owner, post string) (URLs []string, username string, err error) {
	defer browser.CannelAllocator()
	defer browser.CancelTask()
	postURL := fmt.Sprintf("https://vsco.co/%s/media/%s", owner, post)

	var vscoPost VSCOPost

	err = chromedp.Run(browser.Task,
		chromedp.Navigate(postURL),
		chromedp.WaitReady(VSCOScriptSelector),
		chromedp.Evaluate(VSCOScript, &vscoPost),
	)

	if err != nil {
		return URLs, username, err
	}

	media := vscoPost.Medias.ByID[post]
	username = media.PermaSubdomain

	if len(media.VideoURL) > 0 {
		URLs = append(URLs, fmt.Sprintf("https://%s", media.VideoURL))
	} else {
		URLs = append(URLs, fmt.Sprintf("https://%s", media.ResponsiveURL))
	}

	return URLs, username, err
}
