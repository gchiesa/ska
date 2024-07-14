package multipart

import (
	"github.com/gchiesa/ska/internal/part"
	"regexp"
)

// validatePartialsInContent return error if the contentOriginal contains invalid partials
func isValidContent(content []byte) bool {
	reStart := regexp.MustCompile(part.DelimiterStart)
	reEnd := regexp.MustCompile(part.DelimiterEnd)
	reSection := regexp.MustCompile("(?s)" + part.DelimiterStart + "(.*?)" + part.DelimiterEnd) //nolint:goconst //keeping for better readability

	// is valid if the number of section matches with the number of starts and ends
	startEntries := reStart.FindAll(content, -1)
	endEntries := reEnd.FindAll(content, -1)
	sections := reSection.FindAll(content, -1)
	return len(startEntries) == len(sections) && len(endEntries) == len(sections)
}
