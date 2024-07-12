package part

import (
	"encoding/base32"
	"errors"
	"fmt"
	"path/filepath"
)

const DelimiterId = "ska"
const DelimiterStart = "ska-start"
const DelimiterEnd = "ska-end"

type Part struct {
	parentRefFileUri string
	refFilePath      string
	refFileBasename  string
	id               string
	content          []byte
}

var (
	MultipartError = errors.New("error creating multipart")
	InvalidContent = errors.New("invalid content for multipart")
)

//
// Part is the representation of the smallest unit of contentOriginal that is supported by Swanson
// it is delimited by well known placeholders and it will look like the example below
//
// ```
// This is an example
// file.
//
// # swanson-start
// this is a managed partial
// of
// 3 lines
// # swanson-end
//
// this is an unmanaged part
//
// # swanson-start
// this is a managed partial of 1 line
// # swanson-end
//
// this is remaining part
// ```
// in the example there are 2 parts, and they will be parsed and for each partial a new file is created
// that starts with the file name containing the partial and will follow the naming convention below:
// Given the file name is `test-file.txt` the 2 partial will be named:
//
// `test-file.txt.swanson-1`
// `test-file.txt.swanson-2`
//

func NewPart(fromRefFileUri string, id string) *Part {
	idEncoded := base32.StdEncoding.EncodeToString([]byte(id))
	refFileBasename := filepath.Base(fmt.Sprintf("%s.%s-%s", fromRefFileUri, DelimiterId, idEncoded))
	return &Part{
		id:               id,
		parentRefFileUri: fromRefFileUri,
		refFileBasename:  refFileBasename,
	}
}

func (p *Part) RefFileBasename() string {
	return p.refFileBasename
}

func (p *Part) RefFilePath() string {
	return p.refFilePath
}

func (p *Part) Id() string {
	return p.id
}
