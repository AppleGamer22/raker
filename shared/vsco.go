package shared

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

const (
	VSCOMediaSelector      = "img, video"
	VSCOMediaScript        = `document.querySelectorAll("img, video")[1].src`
	VSCOErrorCheckSelector = "p.NotFound-heading"
)

func (raker *Raker) VSCO(owner, post string) (URL string, username string, err error) {
	defer raker.CannelAllocator()
	defer raker.CancelTask()

	timeout, cancel := context.WithTimeout(raker.Task, time.Second*30)
	defer cancel()

	postURL := fmt.Sprintf("https://vsco.co/%s/media/%s", owner, post)
	err = chromedp.Run(timeout,
		chromedp.Navigate(postURL),
		chromedp.WaitNotPresent(VSCOErrorCheckSelector),
		chromedp.WaitReady(VSCOMediaSelector),
		chromedp.Evaluate(VSCOMediaScript, &URL),
		chromedp.TextContent("h4 > a", &username),
	)

	if err == nil {
		URL = fmt.Sprintf("https:%s", URL)
	}

	return URL, username, err
}
