package cleaner

import (
	"regexp"
	"runtime"
	"strings"
)

func Path(input string) string {
	if runtime.GOOS == "windows" {
		input = strings.ReplaceAll(input, `..\`, "")
		input = strings.ReplaceAll(input, `..`, "")
		input = regexp.MustCompile(`[A-Z]:`).ReplaceAllString(input, "")
	} else {
		input = strings.ReplaceAll(input, "../", "")
		input = strings.ReplaceAll(input, "..", "")
	}
	return Line(input)
}
