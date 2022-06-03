package shared

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

const (
	VSCOScriptSelector     = "body > script:nth-child(3)"
	VSCOErrorCheckSelector = "p.NotFound-heading"
)

var VSCOScript = fmt.Sprintf(
	`JSON.parse(document.querySelector("body > script:nth-child(3)").text.slice(%d))`,
	len("window.__PRELOADED_STATE__ = "),
)

type VSCOPost struct {
	Medias struct {
		ByID map[string]struct {
			PermaSubdomain string `json:"permaSubdomain"`
			ResponsiveURL  string `json:"responsiveUrl"`
			VideoURL       string `json:"videoUrl"`
		} `json:"byId"`
	} `json:"medias"`
}

func (raker *Raker) VSCO(owner, post string) (URL string, username string, err error) {
	defer raker.CannelAllocator()
	defer raker.CancelTask()

	timeout, cancel := context.WithTimeout(raker.Task, time.Second*5)
	defer cancel()

	postURL := fmt.Sprintf("https://vsco.co/%s/media/%s", owner, post)
	if err = chromedp.Run(timeout, chromedp.Navigate(postURL)); err != nil {
		return URL, username, err
	}

	timeout, cancel = context.WithTimeout(raker.Task, time.Second*10)
	defer cancel()

	var vscoPost VSCOPost

	err = chromedp.Run(timeout,
		chromedp.WaitNotPresent(VSCOErrorCheckSelector),
		chromedp.WaitReady(VSCOScriptSelector),
		chromedp.Evaluate(VSCOScript, &vscoPost),
	)

	if err != nil {
		return URL, username, err
	}

	media := vscoPost.Medias.ByID[post]
	username = media.PermaSubdomain

	if len(media.VideoURL) > 0 {
		URL = fmt.Sprintf("https://%s", media.VideoURL)
	} else {
		URL = fmt.Sprintf("https://%s", media.ResponsiveURL)
	}

	return URL, username, err
}
