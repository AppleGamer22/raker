package shared

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

const (
	TikTokMediaScript    = `document.querySelector("video").src`
	TikTokUsernameScript = `document.querySelector("h3").innerText`
)

func (raker *Raker) TikTok(owner, post string) (URL string, username string, err error) {
	defer raker.CannelAllocator()
	defer raker.CancelTask()

	timeout, cancel := context.WithTimeout(raker.Task, time.Second*30)
	defer cancel()

	postURL := fmt.Sprintf("https://www.tiktok.com/@%s/video/%s", owner, post)

	err = chromedp.Run(timeout,
		chromedp.Navigate(postURL),
		chromedp.WaitNotPresent("div.error-page"),
		chromedp.WaitReady("video"),
		chromedp.Evaluate(TikTokMediaScript, &URL),
		chromedp.Evaluate(TikTokUsernameScript, &username),
	)

	return URL, username, err
}
