package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var jsonFixture1 = `
{
        "items": [
            {"name": "John Doe", "age": 30},
            {"name": "Jane Doe", "age": 25},
            {"name": "Alice", "age": 27}
        ]
    }
`

var jsonFixture2 = `
{
  "Config" : {
    "BlueprintURI" : "file:///Users/gchiesa/git/swanson/ska-example-template/gotemplate",
    "IgnorePaths" : [ "idea/*", "test-file-to-be-ignored-example.txt", "*.ignored" ]
  },
  "State" : {
    "LastUpdate" : "2024-10-19 14:59:55 +0000 UTC",
    "Variables" : {
      "appName" : "test-app",
      "listOfItems" : "one,two,three",
      "newVersion" : "1.2.3",
      "testFileName" : "example"
    }
  }
}
`

func TestQueryJSONString(t *testing.T) {
	testCases := []struct {
		name        string
		fixtureData *string
		query       string
		result      string
	}{
		{"find names on sample fixture1", &jsonFixture1, "{.items[0].name}", "John Doe"},
		{"find LastUpdate on sample fixture2", &jsonFixture2, "{.State.LastUpdate}", "2024-10-19 14:59:55 +0000 UTC"},
	}
	for _, tc := range testCases {
		result, err := QueryJSONString(*tc.fixtureData, tc.query)
		assert.NoError(t, err)
		assert.Equal(t, tc.result, result)
	}
}
