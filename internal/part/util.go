package part

import (
	"fmt"
	"regexp"
)

const DelimiterID = "ska"
const DelimiterStart = "ska-start"
const DelimiterEnd = "ska-end"

// getPartialsRegexp compiles and returns a regular expression for matching patterns with specific delimiters in a text.
func getPartialsRegexp() *regexp.Regexp {
	// `(?m)(?s)^\s*.{1}\s*%s:(.*?)\s*\n(.*?)^\s*.{1}\s*%s`gm
	pattern := fmt.Sprintf(`(?m)(?s)^\s*.{1}\s*%s:(.*?)\s*\n(.*?)^\s*.{1}\s*%s`, regexp.QuoteMeta(DelimiterStart), regexp.QuoteMeta(DelimiterEnd))
	return regexp.MustCompile(pattern)
}

func buildReplaceRegexp(partialKey string) *regexp.Regexp {
	keyPart := "(.*?)"
	if partialKey != "" {
		keyPart = regexp.QuoteMeta(partialKey)
	}
	pattern := fmt.Sprintf(`(?m)(?s)`+
		`(\s*.{1}\s*%s)`+`:(%s)\s*`+
		`(.*?)`+
		`(^\s*.{1}\s*%s)`, regexp.QuoteMeta(DelimiterStart), keyPart, regexp.QuoteMeta(DelimiterEnd))
	return regexp.MustCompile(pattern)
}
