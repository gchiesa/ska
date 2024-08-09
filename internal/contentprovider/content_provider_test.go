package contentprovider

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestByURI(t *testing.T) {

	var testCases = []struct {
		name         string
		uri          string
		expectedType RemoteContentProvider
	}{
		{
			name:         "github",
			uri:          "https://github.com/gchiesa/test",
			expectedType: &GitHub{},
		},
		{
			name:         "gitlab",
			uri:          "https://gitlab.com/gchiesa/test",
			expectedType: &GitLab{},
		},
		{
			name:         "local",
			uri:          "file:///home/gchiesa/test",
			expectedType: &LocalPath{},
		},
	}
	for _, tc := range testCases {
		cp, err := ByURI(tc.uri)
		assert.Nil(t, err)
		cpType := reflect.TypeOf(cp).String()
		expType := reflect.TypeOf(tc.expectedType).String()
		assert.Equalf(t, expType, cpType, "error in %s test, expected %s, got %s", tc.name, expType, cpType)
	}
}
