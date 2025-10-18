package part

import (
	"encoding/base32"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Part struct {
	parentRefFileURI string
	refFilePath      string
	refFileBasename  string
	id               string
	adoptType        string // e.g. replace-match, inject-before, inject-after
	adoptArg         string // e.g. @start, @end, /^.*$/
	content          []byte
	contentWrapped   []byte
}

var (
	ErrMultipartError = errors.New("error creating multipart")
	ErrInvalidContent = errors.New("invalid content for multipart")
)

// Part is the representation of the smallest unit from the original content
// it is delimited by well-known placeholders, and it will look like the example below
//
// ```
// This is an example
// file.
//
// # ska-start:identifier
// this is a managed partial
// of
// 3 lines
// # ska-end
//
// this is an unmanaged part
//
// # ska-start:identifier2
// this is a managed partial of 1 line
// # ska-end
//
// this is remaining part
// ```
// in the example there are 2 parts, and they will be parsed and for each partial a new file is created
// that starts with the file id containing the partial and will follow the naming convention below:
// Given the file id is `test-file.txt` the 2 partial will be named:
//
// `test-file.txt.ska-1`
// `test-file.txt.ska-2`

func NewPart(fromRefFileURI, id string) *Part {
	idEncoded := base32.StdEncoding.EncodeToString([]byte(id))
	refFileBasename := filepath.Base(fmt.Sprintf("%s.%s-%s", fromRefFileURI, DelimiterID, idEncoded))
	return &Part{
		id:               id,
		parentRefFileURI: fromRefFileURI,
		refFileBasename:  refFileBasename,
	}
}

func (p *Part) SetAdopt(opType, arg string) *Part {
	p.adoptType = opType
	p.adoptArg = arg
	return p
}

func (p *Part) RefFileBasename() string {
	return p.refFileBasename
}

func (p *Part) RefFilePath() string {
	return p.refFilePath
}

func (p *Part) ID() string {
	return p.id
}

func (p *Part) AdoptType() string { return p.adoptType }
func (p *Part) AdoptArg() string  { return p.adoptArg }

func (p *Part) WithContent(content []byte) *Part {
	p.content = content
	return p
}

func (p *Part) SetFilePath(path string) *Part {
	p.refFilePath = path
	return p
}

func (p *Part) CreateFile() error {
	return os.WriteFile(p.refFilePath, p.content, 0o644)
}

// ParseParts extracts structured parts from the originalContent based on specific delimiters and directives in the headers.
// It returns a slice of Part with metadata and content, or an error if parsing fails.
func ParseParts(originalContent []byte, parentRefFileURI string) ([]Part, error) {
	var parts []Part
	reSection := getPartialsRegexp()

	matches := reSection.FindAllSubmatch(originalContent, -1)

	for _, match := range matches {
		originalHeader := strings.TrimSpace(string(match[1]))
		content := string(match[2])
		sectionName, adoptType, adoptArg, err := parseHeader(originalHeader)
		if err != nil {
			return nil, err
		}
		partial := NewPart(parentRefFileURI, sectionName).WithContent([]byte(content))
		if adoptType != "" {
			partial = partial.SetAdopt(adoptType, adoptArg)
		}
		parts = append(parts, *partial)
	}
	return parts, nil
}

// parseHeader splits the part header "<section> [+ directive:arg]" into components
func parseHeader(header string) (sectionName, adoptType, adoptArg string, err error) {
	headerToParse := strings.TrimSpace(header)
	headerToParse = strings.TrimPrefix(headerToParse, ":")

	adoptType, adoptArg = "", ""
	if strings.Contains(headerToParse, "+") {
		parts := strings.SplitN(headerToParse, "+", 2)
		sectionName = strings.TrimSpace(parts[0])
		directive := strings.TrimSpace(parts[1])
		// normalize
		directive = strings.TrimSpace(directive)
		// supported directives
		supported := []string{"ska-inject-after:", "ska-inject-before:", "ska-replace-match:"}
		for _, pref := range supported {
			if strings.HasPrefix(directive, pref) {
				adoptType = strings.TrimSuffix(pref, ":")
				adoptArg = strings.TrimSpace(strings.TrimPrefix(directive, pref))
				break
			}
		}
	} else {
		sectionName = headerToParse
	}
	if sectionName == "" {
		err = ErrInvalidContent
		return
	}
	return
}
