package cleaner

import "regexp"

func MongoQuery(input string) string {
	input = regexp.MustCompile(`^\$|\.`).ReplaceAllString(input, "")
	return Line(input)
}
