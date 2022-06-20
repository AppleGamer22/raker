package cleaner

import "strings"

func Line(entry string) string {
	escapedEntry := strings.ReplaceAll(entry, "\n", "")
	escapedEntry = strings.ReplaceAll(escapedEntry, "\r", "")
	return escapedEntry
}
