package yamlmerge

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestMergeYAML_scalarOverride(t *testing.T) {
	t.Parallel()
	dst := []byte(`test:
  key: abc
  version: base
`)
	src := []byte(`test:
  version: updated
`)
	merged, paths, err := MergeYAML(dst, src)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, yaml.Unmarshal(merged, &result))

	testMap := result["test"].(map[string]interface{})
	assert.Equal(t, "abc", testMap["key"])         // preserved
	assert.Equal(t, "updated", testMap["version"]) // overridden

	assert.Contains(t, paths, "test.version")
	assert.NotContains(t, paths, "test.key")
}

func TestMergeYAML_preservesExtraKeys(t *testing.T) {
	t.Parallel()
	dst := []byte(`test2:
    key2: value2
userAdded:
    key3: value3
`)
	src := []byte(`test2:
    key2: newValue
`)
	merged, paths, err := MergeYAML(dst, src)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, yaml.Unmarshal(merged, &result))

	test2 := result["test2"].(map[string]interface{})
	assert.Equal(t, "newValue", test2["key2"]) // overridden

	userAdded, ok := result["userAdded"]
	assert.True(t, ok, "userAdded key should be preserved")
	assert.Equal(t, "value3", userAdded.(map[string]interface{})["key3"])

	assert.Contains(t, paths, "test2.key2")
}

func TestMergeYAML_withLeadingIndent(t *testing.T) {
	t.Parallel()
	// Simulates content extracted from inside a managed block (4-space common indent)
	dst := []byte("    test2:\n        key2: value2\n    userAdded:\n        key3: value3\n")
	src := []byte("    test2:\n        key2: rendered\n")

	merged, paths, err := MergeYAML(dst, src)
	require.NoError(t, err)

	// The merged content should still have the same indentation level.
	assert.Contains(t, string(merged), "    ")
	assert.Contains(t, paths, "test2.key2")

	// Parse after stripping indent to verify correctness.
	stripped, _ := stripCommonIndent(merged)
	var result map[string]interface{}
	require.NoError(t, yaml.Unmarshal(stripped, &result))

	assert.Equal(t, "rendered", result["test2"].(map[string]interface{})["key2"])
	assert.Equal(t, "value3", result["userAdded"].(map[string]interface{})["key3"])
}

func TestMergeYAML_noChanges(t *testing.T) {
	t.Parallel()
	dst := []byte(`key: value
`)
	src := []byte(`key: value
`)
	_, paths, err := MergeYAML(dst, src)
	require.NoError(t, err)
	assert.Empty(t, paths)
}

func TestMergeYAML_emptySrc(t *testing.T) {
	t.Parallel()
	dst := []byte(`key: value
`)
	merged, paths, err := MergeYAML(dst, []byte("   \n"))
	require.NoError(t, err)
	assert.Equal(t, dst, merged)
	assert.Empty(t, paths)
}

func TestMergeYAML_addNewKey(t *testing.T) {
	t.Parallel()
	dst := []byte(`existing: yes
`)
	src := []byte(`newKey: added
`)
	merged, paths, err := MergeYAML(dst, src)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, yaml.Unmarshal(merged, &result))
	assert.Equal(t, "yes", result["existing"])
	assert.Equal(t, "added", result["newKey"])
	assert.Contains(t, paths, "newKey")
}

func TestStripCommonIndent(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name           string
		input          string
		expectedIndent int
		expectedLine0  string
	}{
		{"no indent", "key: val\n", 0, "key: val"},
		{"2-space indent", "  key: val\n  other: x\n", 2, "key: val"},
		{"4-space indent", "    key: val\n    other: x\n", 4, "key: val"},
		{"mixed: preserves relative", "    a: 1\n      b: 2\n", 4, "a: 1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stripped, indent := stripCommonIndent([]byte(tc.input))
			assert.Equal(t, tc.expectedIndent, indent)
			lines := splitLines(string(stripped))
			if len(lines) > 0 {
				assert.Equal(t, tc.expectedLine0, lines[0])
			}
		})
	}
}

func splitLines(s string) []string {
	var out []string
	for _, l := range []string(splitStr(s)) {
		if l != "" {
			out = append(out, l)
		}
	}
	return out
}

func splitStr(s string) []string {
	result := []string{}
	cur := ""
	for _, c := range s {
		if c == '\n' {
			result = append(result, cur)
			cur = ""
		} else {
			cur += string(c)
		}
	}
	if cur != "" {
		result = append(result, cur)
	}
	return result
}
