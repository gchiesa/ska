package multipart

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const fileExampleGood01 = `
This is an example 
file.

# ska-start
this is a managed partial
of 
3 lines
# ska-end

this is an unmanaged part

# ska-start
this is a managed partial of 1 line
# ska-end 

this is remaining part
`

const fileExampleGood02 = `
[test]
This is an example of init file

[section]
; ska-start 
this is a managed partial
; ska-end
`

const fileExampleWrong01 = `
This is an example 
file.

# ska-start:key01
this is a managed partial
of 
3 lines

# ska-start:key02
this another managed partial of 1 line
# ska-end 

this is remaining part
`

const fileExampleWrong02 = `
This is an example 
file.

# ska-start
this another managed partial of 1 line with no key
# ska-end 

this is remaining part
`

const fileExampleWrong03 = `
This is an example 
file.

# ska-start:key01
this another managed partial of 1 line with no key
# ska-end 

# ska-start:key01
this is remaining part with duplicate key
# ska-end
`

func TestMultipartValidation(t *testing.T) {
	testCases := []struct {
		name        string
		content     []byte
		expectValid bool
	}{
		{
			name:        "fileExampleGood01",
			content:     []byte(fileExampleGood01),
			expectValid: true,
		},
		{
			name:        "fileExampleGood02",
			content:     []byte(fileExampleGood01),
			expectValid: true,
		},
		{
			name:        "fileExampleWrong01",
			content:     []byte(fileExampleWrong01),
			expectValid: false,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectValid, isValidContent(tc.content))
	}
}
