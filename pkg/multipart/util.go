package multipart

import (
	"fmt"
	"github.com/gchiesa/ska/pkg/part"
	"regexp"
)

func getPartialsRegexp() *regexp.Regexp {
	// `(?m)(?s)^\s*.{1}\s*%s:(.*?)\s*\n(.*?)^\s*.{1}\s*%s`gm
	pattern := fmt.Sprintf(`(?m)(?s)^\s*.{1}\s*%s:(.*?)\s*\n(.*?)^\s*.{1}\s*%s`, regexp.QuoteMeta(part.DelimiterStart), regexp.QuoteMeta(part.DelimiterEnd))
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
		`(^\s*.{1}\s*%s)`, regexp.QuoteMeta(part.DelimiterStart), keyPart, regexp.QuoteMeta(part.DelimiterEnd))
	return regexp.MustCompile(pattern)
}
