package contentprovider

import (
	"github.com/apex/log"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/stretchr/testify/assert"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
	"testing"
)

const (
	GitHubTestRepository = "https://github.com/gchiesa/ska-example-template@main"
)

func TestGitHubDownloadContent(t *testing.T) {
	r, err := recorder.New("fixtures/content-provider-github")
	if err != nil {
		log.Fatalf("error creating recorder: %v", err)
	}
	defer func(r *recorder.Recorder) { _ = r.Stop() }(r)

	// fake httpClient
	httpClient := r.GetDefaultClient()

	// use the mock client
	client.InstallProtocol("https", githttp.NewClient(httpClient))

	gh, err := NewGitHub(GitHubTestRepository)
	if err != nil {
		log.Fatalf("error creating GitHub client: %v", err)
	}
	err = gh.DownloadContent()
	assert.NoError(t, err)
}

func TestGithub_validateRemoteURI(t *testing.T) {
	testcases := []struct {
		name              string
		uri               string
		expectedTupleFunc func(t *testing.T, cp *GitHub)
		expectedErrFunc   func(t *testing.T, err error)
	}{
		{"Given URI without branch should return an error",
			"https://github.com/gchiesa/test",
			func(t *testing.T, cp *GitHub) { return },
			func(t *testing.T, err error) { assert.Error(t, err) },
		},
		{
			"Given URI with branch should return the tuple and no errors with filePath empty",
			"https://github.com/gchiesa/test@branch123",
			func(t *testing.T, cp *GitHub) {
				assert.NotEmpty(t, cp.remoteURI)
				assert.Equal(t, "branch123", cp.repositoryRef)
				assert.Empty(t, cp.repositoryFilePath)
			},
			func(t *testing.T, err error) { assert.NoError(t, err) },
		},
		{
			"Given URI with filepath and reference, the tuple should match the parts and no errors",
			"https://github.com/gchiesa/test//path_from_root/file1@main",
			func(t *testing.T, cp *GitHub) {
				assert.NotEmpty(t, cp.remoteURI)
				assert.NotEmpty(t, cp.repositoryRef)
				assert.Equal(t, "path_from_root/file1", cp.repositoryFilePath)
			},
			func(t *testing.T, err error) { assert.NoError(t, err) },
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cp, err := NewGitHub(tc.uri)
			assert.NoError(t, err)
			err = cp.validateRemoteURI(tc.uri)
			tc.expectedTupleFunc(t, cp)
			tc.expectedErrFunc(t, err)
		})
	}

}
