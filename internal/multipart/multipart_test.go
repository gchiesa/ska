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
	t.Parallel()
	// Table-driven tests (single scenario here, structured for extension)
	cases := []struct {
		name           string
		blueprintFile  string
		destBeforeFile string
	}{
		{
			name:           "inject-after-start-idempotent",
			blueprintFile:  filepath.Join("fixtures", "adopt-inject-after-start-blueprint.txt"),
			destBeforeFile: filepath.Join("fixtures", "adopt-inject-after-start-dest-before.txt"),
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tmp := t.TempDir()
			// load fixtures
			bpContent, err := os.ReadFile(tc.blueprintFile)
			if err != nil {
				t.Fatalf("read blueprint: %v", err)
			}
			beforeContent, err := os.ReadFile(tc.destBeforeFile)
			if err != nil {
				t.Fatalf("read before: %v", err)
			}

			// write blueprint into temp area so ref files are created there
			bpPath := writeTempFile(t, tmp, "blue.txt", string(bpContent))
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
			dest := writeTempFile(t, tmp, "dest.txt", string(beforeContent))
			// compile to destination (will read dest content as original)
			if err := mpart.CompileToFile(dest, false); err != nil {
				t.Fatalf("compile: %v", err)
			}
			data, _ := os.ReadFile(dest)
			got := string(data)
			// ensure header has no directive (idempotency)
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
		})
	}
}

func TestAdoptReplaceMatchWholeAndGroup(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name           string
		blueprintFile  string
		destBeforeFile string
		partialContent string
	}
	cases := []testCase{
		{
			name:           "replace-whole-match",
			blueprintFile:  filepath.Join("fixtures", "replace-whole-blueprint.txt"),
			destBeforeFile: filepath.Join("fixtures", "replace-whole-dest-before.txt"),
			partialContent: "X\n",
		},
		{
			name:           "replace-group-match",
			blueprintFile:  filepath.Join("fixtures", "replace-group-blueprint.txt"),
			destBeforeFile: filepath.Join("fixtures", "replace-group-dest-before.txt"),
			partialContent: "Y\n",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tmp := t.TempDir()
			bpContent, err := os.ReadFile(tc.blueprintFile)
			if err != nil {
				t.Fatalf("read blueprint: %v", err)
			}
			beforeContent, err := os.ReadFile(tc.destBeforeFile)
			if err != nil {
				t.Fatalf("read before: %v", err)
			}

			bpPath := writeTempFile(t, tmp, "blue.txt", string(bpContent))
			m, err := NewMultipartFromFile(bpPath, "blue.txt")
			if err != nil {
				t.Fatalf("multipart: %v", err)
			}
			if err := m.ParseParts(); err != nil {
				t.Fatalf("parse: %v", err)
			}
			for _, p := range m.Parts() {
				writeTempFile(t, filepath.Dir(bpPath), p.RefFileBasename(), tc.partialContent)
			}
			dest := writeTempFile(t, tmp, "dest.txt", string(beforeContent))
			if err := m.CompileToFile(dest, false); err != nil {
				t.Fatalf("compile: %v", err)
			}
			gotB, _ := os.ReadFile(dest)
			got := string(gotB)
			switch tc.name {
			case "replace-whole-match":
				assert.Contains(t, got, "# ska-start:r1\nX\n# ska-end\n")
				assert.NotContains(t, got, "PLACEHOLDER")
			case "replace-group-match":
				assert.Contains(t, got, "# ska-start:r2\nY\n# ska-end\n")
				assert.NotContains(t, got, "target")
			}
		})
	}
}

func TestReplaceMatchMultiline(t *testing.T) {
	// Test that replaceMatch correctly handles multiline content with ^ and $ anchors
	base := []byte(`FROM alpine:3.20
COPY my-test-app /usr/bin/my-test-app
ENTRYPOINT ["/usr/bin/my-test-app"]
`)
	regex := `^FROM .*$`
	payload := `# ska-start:new-base
FROM ubuntu:latest
# ska-end`

	result := replaceMatch(base, regex, payload)

	expected := `# ska-start:new-base
FROM ubuntu:latest
# ska-end
COPY my-test-app /usr/bin/my-test-app
ENTRYPOINT ["/usr/bin/my-test-app"]
`
	assert.Equal(t, expected, string(result))
}
