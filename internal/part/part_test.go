package part

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHeader(t *testing.T) {
	t.Parallel()
	type Expected struct {
		sectionName         string
		adoptType, adoptArg string
	}

	testCases := []struct {
		name        string
		header      string
		expected    Expected
		expectedErr bool
	}{
		{"no section id - should return error", "", Expected{sectionName: "", adoptType: "", adoptArg: ""}, true},
		{"no section id with adopt meta - should return error", ": + ska-inject-before:@start", Expected{sectionName: "", adoptType: "", adoptArg: ""}, true},
		{"only id", ":my-section", Expected{sectionName: "my-section", adoptType: "", adoptArg: ""}, false},
		{"id with adopt meta", ":my-section + ska-inject-before:@start", Expected{sectionName: "my-section", adoptType: "ska-inject-before", adoptArg: "@start"}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sectionName, adoptType, adoptArg, err := parseHeader(tc.header)
			if tc.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected.sectionName, sectionName)
			assert.Equal(t, tc.expected.adoptType, adoptType)
			assert.Equal(t, tc.expected.adoptArg, adoptArg)
		})
	}
}

func TestParseParts(t *testing.T) {
	t.Parallel()
	type ExpectedPart struct {
		id        string
		adoptType string
		adoptArg  string
	}
	testCases := []struct {
		name          string
		fixtureFile   string
		expectedParts []ExpectedPart
	}{
		{
			name:          "no parts found when no colon after start",
			fixtureFile:   filepath.Join("fixtures", "header-default.txt"),
			expectedParts: []ExpectedPart{},
		},
		{
			name:          "single named part",
			fixtureFile:   filepath.Join("fixtures", "header-with-named-part.txt"),
			expectedParts: []ExpectedPart{{id: "my-part", adoptType: "", adoptArg: ""}},
		},
		{
			name:          "single named part with adopt",
			fixtureFile:   filepath.Join("fixtures", "header-with-named-part-and-action.txt"),
			expectedParts: []ExpectedPart{{id: "my-part", adoptType: "ska-inject-before", adoptArg: "@start"}},
		},
		{
			name:          "multiple parts in single file",
			fixtureFile:   filepath.Join("fixtures", "multi-parts.txt"),
			expectedParts: []ExpectedPart{{id: "part-1"}, {id: "part-2"}},
		},
		{
			name:          "multiple parts with adopt directives",
			fixtureFile:   filepath.Join("fixtures", "multi-parts-with-adopt.txt"),
			expectedParts: []ExpectedPart{{id: "alpha", adoptType: "ska-inject-before", adoptArg: "@start"}, {id: "beta", adoptType: "ska-inject-after", adoptArg: "@end"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content, err := os.ReadFile(tc.fixtureFile)
			assert.NoError(t, err)
			parts, err := ParseParts(content, tc.fixtureFile)
			assert.NoError(t, err)
			assert.Equal(t, len(tc.expectedParts), len(parts))
			for i, ep := range tc.expectedParts {
				assert.Equal(t, ep.id, parts[i].ID())
				assert.Equal(t, ep.adoptType, parts[i].AdoptType())
				assert.Equal(t, ep.adoptArg, parts[i].AdoptArg())
			}
		})
	}
}
