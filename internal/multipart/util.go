package multipart

import (
	"fmt"
	"regexp"

	"github.com/gchiesa/ska/internal/part"
)

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
