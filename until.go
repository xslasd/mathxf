package mathxf

import "regexp"

func containsAtLeastOneLetter(s string) bool {
	pattern := "[a-zA-Z0-9_]*[a-zA-Z][a-zA-Z0-9_]*"
	matcher := regexp.MustCompile(pattern)
	return matcher.MatchString(s)
}
