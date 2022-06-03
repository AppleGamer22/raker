package shared

import (
	"context"

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

func NewRaker(execPath, userDateDir string, debug, incognito bool) (Raker, error) {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.UserDataDir(userDateDir),
		chromedp.ExecPath(execPath),
		chromedp.Flag("incognito", incognito),
		chromedp.Flag("headless", !debug),
	)

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
