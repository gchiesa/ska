package multipart

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func writeTempFile(t *testing.T, dir, name, content string) string {
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return p
}

func TestAdoptInjectAfterStartAndIdempotency(t *testing.T) {
	tmp := t.TempDir()
	// blueprint file with adopt directive
	blueprint := `# ska-start:my-section + ska-inject-after:<@start>
 managed v1
 # ska-end
 `
	bpPath := writeTempFile(t, tmp, "blue.txt", blueprint)
	mpart, err := NewMultipartFromFile(bpPath, "blue.txt")
	if err != nil {
		t.Fatalf("multipart: %v", err)
	}
	if err := mpart.ParseParts(); err != nil {
		t.Fatalf("parse: %v", err)
	}
	if assert.Greater(t, len(mpart.Parts()), 0, "no parts parsed") {
		assert.Equal(t, "ska-inject-after", mpart.Parts()[0].AdoptType())
		assert.Equal(t, "@start", mpart.Parts()[0].AdoptArg())
	}
	// create compiled partial file content (what renderer would create)
	for _, p := range mpart.Parts() {
		writeTempFile(t, filepath.Dir(bpPath), p.RefFileBasename(), "managed v1\n")
	}
	// destination file initially without managed block
	dest := writeTempFile(t, tmp, "dest.txt", "Header\nBody\nFooter\n")
	// compile to destination (will read dest content as original)
	if err := mpart.CompileToFile(dest, false); err != nil {
		t.Fatalf("compile: %v", err)
	}
	data, _ := os.ReadFile(dest)
	got := string(data)
	// expect injection at start (top of file)
	assert.Contains(t, got, "# ska-start:my-section\nmanaged v1\n# ska-end\n")
	assert.True(t, strings.HasPrefix(got, "# ska-start:my-section\nmanaged v1\n# ska-end\n"))
	// ensure header has no directive (idempotency)
	assert.NotContains(t, got, "+ ska-inject-after")

	// Update partial content and run again: should replace existing block
	for _, p := range mpart.Parts() {
		writeTempFile(t, filepath.Dir(bpPath), p.RefFileBasename(), "managed v2\n")
	}
	if err := mpart.CompileToFile(dest, false); err != nil {
		t.Fatalf("compile2: %v", err)
	}
	data2, _ := os.ReadFile(dest)
	got2 := string(data2)
	assert.Contains(t, got2, "# ska-start:my-section\nmanaged v2\n# ska-end\n")
	assert.NotContains(t, got2, "managed v1\n")
}

func TestAdoptReplaceMatchWholeAndGroup(t *testing.T) {
	tmp := t.TempDir()
	// Whole match replacement
	blueprint1 := `# ska-start:r1 + ska-replace-match:<PLACEHOLDER>
 X
 # ska-end
 `
	bp1 := writeTempFile(t, tmp, "b1.txt", blueprint1)
	m1, _ := NewMultipartFromFile(bp1, "b1.txt")
	_ = m1.ParseParts()
	for _, p := range m1.Parts() {
		writeTempFile(t, filepath.Dir(bp1), p.RefFileBasename(), "X\n")
	}
	d1 := writeTempFile(t, tmp, "d1.txt", "pre\nPLACEHOLDER\npost\n")
	_ = m1.CompileToFile(d1, false)
	b1, _ := os.ReadFile(d1)
	got1 := string(b1)
	assert.Contains(t, got1, "# ska-start:r1\nX\n# ska-end\n")
	assert.NotContains(t, got1, "PLACEHOLDER")

	// Single capture group replacement (basic presence check)
	blueprint2 := `# ska-start:r2 + ska-replace-match:<MIDDLE(.*)END>
	Y
	# ska-end
	`
	bp2 := writeTempFile(t, tmp, "b2.txt", blueprint2)
	m2, _ := NewMultipartFromFile(bp2, "b2.txt")
	_ = m2.ParseParts()
	for _, p := range m2.Parts() {
		writeTempFile(t, filepath.Dir(bp2), p.RefFileBasename(), "Y\n")
	}
	d2 := writeTempFile(t, tmp, "d2.txt", "pre\nMIDDLEtargetEND\npost\n")
	_ = m2.CompileToFile(d2, false)
	b2, _ := os.ReadFile(d2)
	got2 := string(b2)
	assert.Contains(t, got2, "# ska-start:r2\nY\n# ska-end\n")
	assert.NotContains(t, got2, "target")
}
