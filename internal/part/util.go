package part

import (
	"fmt"
	"regexp"
)

const DelimiterID = "ska"
const DelimiterStart = "ska-start"
const DelimiterEnd = "ska-end"

// getPartialsRegexp compiles and returns a regular expression for matching patterns with specific delimiters in a text.
// It supports an optional bracket modifier group between the start delimiter and the colon, e.g.:
//
//	`# ska-start[engine:yaml-merge]:identifier`  (new format with modifiers)
//	`# ska-start:identifier`                     (legacy format without modifiers)
//
// Capture groups: 1=modifiers (optional), 2=header (identifier + optional directive), 3=content
func getPartialsRegexp() *regexp.Regexp {
	// `(?m)(?s)^\s*.{1}\s*ska-start(?:\[(.*?)\])?:(.*?)\s*\n(.*?)^\s*.{1}\s*ska-end`
	pattern := fmt.Sprintf(`(?m)(?s)^\s*.{1}\s*%s(?:\[(.*?)\])?:(.*?)\s*\n(.*?)^\s*.{1}\s*%s`,
		regexp.QuoteMeta(DelimiterStart), regexp.QuoteMeta(DelimiterEnd))
	return regexp.MustCompile(pattern)
}
