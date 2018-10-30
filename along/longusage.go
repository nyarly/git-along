package along

import (
	"regexp"
)

var longFix = regexp.MustCompile(`(?m)^[ \t]*`)

func longUsage(s string) string {
	return longFix.ReplaceAllString(s, "")
}
