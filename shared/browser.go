package shared

import (
	"context"
	"os"
	"runtime"

	"github.com/chromedp/chromedp"
)

type Raker struct {
	Debug           bool
	Incognito       bool
	Allocator       context.Context
	CannelAllocator context.CancelFunc
	Task            context.Context
	CancelTask      context.CancelFunc
}

var (
	UserDataDirectory string
	// ExecutablePath    string
)

func FindExecutablePath() string {
	switch runtime.GOOS {
	case "darwin":
		return "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	case "windows":
		return `C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`
	default:
		if _, err := os.Stat("/usr/bin/chromium"); err == nil {
			return "/usr/bin/chromium"
		} else if _, err := os.Stat("/usr/bin/chromium-browser"); err == nil {
			return "/usr/bin/chromium-browser"
		}
		return "/opt/google/chrome/google-chrome"
	}
}

func NewRaker(userDateDir string, debug, incognito bool) (Raker, error) {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("incognito", incognito),
		chromedp.Flag("headless", !debug),
		chromedp.ExecPath(FindExecutablePath()),
	)

	if userDateDir != "" {
		opts = append(opts, chromedp.UserDataDir(userDateDir))
	}

	allocator, cancelAllocator := chromedp.NewExecAllocator(context.Background(), opts...)
	task, cancelTask := chromedp.NewContext(allocator)

	raker := Raker{
		Debug:           debug,
		Incognito:       incognito,
		Allocator:       allocator,
		CannelAllocator: cancelAllocator,
		Task:            task,
		CancelTask:      cancelTask,
	}

	err := chromedp.Run(raker.Task, chromedp.Evaluate("delete Object.getPrototypeOf(navigator).webdriver", nil))
	return raker, err
}
