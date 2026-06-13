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
			mpart, err := NewMultipartFromFile(bpPath, "blue.txt", nil)
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
			if err := mpart.CompileToFile(dest); err != nil {
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
			if err := mpart.CompileToFile(dest); err != nil {
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
			m, err := NewMultipartFromFile(bpPath, "blue.txt", nil)
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
			if err := m.CompileToFile(dest); err != nil {
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

func TestYAMLMergeEngine_fullBlock(t *testing.T) {
	t.Parallel()
	// Blueprint has [engine:yaml-merge] modifier; compiled partial replaces test.version but
	// the destination's extra key (extraKey) should be preserved.
	tmp := t.TempDir()

	bpContent, err := os.ReadFile(filepath.Join("fixtures", "yaml-merge-blueprint.txt"))
	if err != nil {
		t.Fatalf("read blueprint: %v", err)
	}
	beforeContent, err := os.ReadFile(filepath.Join("fixtures", "yaml-merge-dest-before.txt"))
	if err != nil {
		t.Fatalf("read before: %v", err)
	}

	bpPath := writeTempFile(t, tmp, "blue.yaml", string(bpContent))
	m, err := NewMultipartFromFile(bpPath, "blue.yaml", nil)
	if err != nil {
		t.Fatalf("multipart: %v", err)
	}
	if err := m.ParseParts(); err != nil {
		t.Fatalf("parse: %v", err)
	}
	require := assert.New(t) // use assert as require for brevity
	require.Equal(1, len(m.Parts()))
	require.Equal("yaml-merge", m.Parts()[0].Engine())

	// Write compiled partial: the rendered template output (test.version = v2.0).
	partialContent := "test:\n  key: abc\n  version: v2.0\n"
	for _, p := range m.Parts() {
		writeTempFile(t, filepath.Dir(bpPath), p.RefFileBasename(), partialContent)
	}

	dest := writeTempFile(t, tmp, "dest.yaml", string(beforeContent))
	if err := m.CompileToFile(dest); err != nil {
		t.Fatalf("compile: %v", err)
	}

	data, _ := os.ReadFile(dest)
	got := string(data)

	// version should be updated
	assert.Contains(t, got, "version: v2.0")
	// extra key added by user should be preserved
	assert.Contains(t, got, "extraKey: preserved")
	// preamble and postamble survive
	assert.Contains(t, got, "preamble: true")
	assert.Contains(t, got, "postamble: true")
}

func TestYAMLMergeEngine_partialSection(t *testing.T) {
	t.Parallel()
	// Blueprint has a partial section with yaml-merge engine.
	// Destination has userAdded key that must survive the merge.
	tmp := t.TempDir()

	bpContent, err := os.ReadFile(filepath.Join("fixtures", "yaml-merge-partial-blueprint.txt"))
	if err != nil {
		t.Fatalf("read blueprint: %v", err)
	}
	beforeContent, err := os.ReadFile(filepath.Join("fixtures", "yaml-merge-partial-dest-before.txt"))
	if err != nil {
		t.Fatalf("read before: %v", err)
	}

	bpPath := writeTempFile(t, tmp, "blue.yaml", string(bpContent))
	m, err := NewMultipartFromFile(bpPath, "blue.yaml", nil)
	if err != nil {
		t.Fatalf("multipart: %v", err)
	}
	if err := m.ParseParts(); err != nil {
		t.Fatalf("parse: %v", err)
	}
	assert.Equal(t, 1, len(m.Parts()))
	assert.Equal(t, "yaml-merge", m.Parts()[0].Engine())

	// Write compiled partial (rendered template): test2.key2 updated to "rendered".
	partialContent := "    test2:\n        key2: rendered\n"
	for _, p := range m.Parts() {
		writeTempFile(t, filepath.Dir(bpPath), p.RefFileBasename(), partialContent)
	}

	dest := writeTempFile(t, tmp, "dest.yaml", string(beforeContent))
	if err := m.CompileToFile(dest); err != nil {
		t.Fatalf("compile: %v", err)
	}

	data, _ := os.ReadFile(dest)
	got := string(data)

	// key2 should be updated
	assert.Contains(t, got, "key2: rendered")
	// userAdded key must be preserved
	assert.Contains(t, got, "userAdded")
	assert.Contains(t, got, "key3: value3")
	// surrounding YAML structure preserved
	assert.Contains(t, got, "version: base")
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

func TestYAMLMergeEngine_adoptReplaceMatch(t *testing.T) {
	t.Parallel()
	// Destination file has NO managed block (first-time adoption via ska-replace-match).
	// Blueprint uses yaml-merge engine with ska-replace-match to adopt the stack: section.
	// The destination has stack.placement.primary_vpc_cidr that is NOT in the compiled partial.
	// After adoption, primary_vpc_cidr must be preserved (not lost due to plain replacement).
	tmp := t.TempDir()

	// Correct syntax: [engine:yaml-merge] in brackets, + directive in the header after ':'
	blueprintContent := strings.TrimSpace(`
# ska-start[engine:yaml-merge]:stack-section + ska-replace-match:(?s)^stack:.*
stack:
  blueprint_ref: v1.0
  team: platform
# ska-end
`) + "\n"

	// Destination: plain YAML, no ska markers, has extra keys not in blueprint
	destContent := `stack:
  blueprint_ref: old-ref
  team: old-team
  placement:
    primary_vpc_cidr: 10.128.0.0/18
    provider: aws
`

	// Compiled partial: only blueprint keys, no primary_vpc_cidr
	partialContent := `stack:
  blueprint_ref: v1.21.5
  team: aws-platform
`

	bpPath := writeTempFile(t, tmp, "blue.yaml", blueprintContent)
	m, err := NewMultipartFromFile(bpPath, "blue.yaml", nil)
	if err != nil {
		t.Fatalf("multipart: %v", err)
	}
	if err := m.ParseParts(); err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(m.Parts()) == 0 {
		t.Fatal("expected at least one part")
	}
	assert.Equal(t, "yaml-merge", m.Parts()[0].Engine())
	assert.Equal(t, "ska-replace-match", m.Parts()[0].AdoptType())

	for _, p := range m.Parts() {
		writeTempFile(t, filepath.Dir(bpPath), p.RefFileBasename(), partialContent)
	}

	dest := writeTempFile(t, tmp, "dest.yaml", destContent)
	if err := m.CompileToFile(dest); err != nil {
		t.Fatalf("compile: %v", err)
	}

	data, _ := os.ReadFile(dest)
	got := string(data)
	t.Logf("Result:\n%s", got)

	// Blueprint keys updated
	assert.Contains(t, got, "blueprint_ref: v1.21.5", "blueprint_ref should be updated from compiled partial")
	assert.Contains(t, got, "team: aws-platform", "team should be updated from compiled partial")
	// Dst-only keys preserved
	assert.Contains(t, got, "primary_vpc_cidr: 10.128.0.0/18", "primary_vpc_cidr (dst-only) must be preserved after adoption")
	assert.Contains(t, got, "provider: aws", "provider (dst-only) must be preserved after adoption")
	// Result is wrapped in a managed block
	assert.Contains(t, got, "# ska-start:stack-section", "result must be wrapped in a managed block")
	assert.Contains(t, got, "# ska-end", "result must be wrapped in a managed block")
}
