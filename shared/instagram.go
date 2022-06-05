package shared

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

const (
	InstagramErrorCheckSelector = "div.error-container"
	InstagramMediaSelector      = "article img[srcset], video"
	InstagramMediaScript        = `Array.from(document.querySelectorAll("article img[srcset], video")).map(element => element.src)`
	InstagramUsernameScript     = `document.querySelectorAll("a")[1].innerText`
)

func (raker *Raker) Instagram(post string) (URLs []string, username string, err error) {
	defer raker.CannelAllocator()
	defer raker.CancelTask()

	timeout, cancel := context.WithTimeout(raker.Task, time.Second*30)
	defer cancel()

	postURL := fmt.Sprintf("https://www.instagram.com/p/%s", post)

	err = chromedp.Run(timeout,
		chromedp.Navigate(postURL),
		chromedp.WaitNotPresent(InstagramErrorCheckSelector),
		chromedp.WaitReady(InstagramMediaSelector),
		chromedp.Evaluate(InstagramMediaScript, &URLs),
		chromedp.Evaluate(InstagramUsernameScript, &username),
	)

	return URLs, username, err
}

func (raker *Raker) InstagramSignIn(username, password string) error {
	defer raker.CannelAllocator()
	defer raker.CancelTask()

	timeout, cancel := context.WithTimeout(raker.Task, time.Second*30)
	defer cancel()

	return chromedp.Run(timeout,
		chromedp.Navigate("https://www.instagram.com/accounts/login/"),
		chromedp.WaitVisible(`input[name="username"]`),
		chromedp.SendKeys(`input[name="username"]`, username),
		chromedp.SendKeys(`input[name="password"]`, password),
		chromedp.Click(`button[type="submit"]`),
		chromedp.WaitVisible("button.sqdOP"),
		chromedp.Click("button.sqdOP"),
		chromedp.WaitVisible(fmt.Sprintf(`a:contains("%s")`, username)),
	)
}

func (raker *Raker) InstagramSignOut(username string) error {
	defer raker.CannelAllocator()
	defer raker.CancelTask()

	timeout, cancel := context.WithTimeout(raker.Task, time.Second*30)
	defer cancel()

	profileURL := fmt.Sprintf("https://www.instagram.com/%s", username)

	return chromedp.Run(timeout,
		chromedp.Navigate(profileURL),
		chromedp.WaitVisible(`svg[aria-label="Options"]`),
		chromedp.Click(`svg[aria-label="Options"]`),
		chromedp.WaitVisible("button:nth-child(9)"),
		chromedp.Click("button:nth-child(9)"),
		chromedp.WaitVisible(`button[type="submit"], button[type="button"]`),
	)
}
