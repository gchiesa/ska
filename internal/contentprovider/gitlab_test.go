package contentprovider

import (
	"github.com/apex/log"
	"github.com/stretchr/testify/assert"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
	"testing"
)

const (
	GitLabTestRepository = "https://gitlab.com/gchiesa/test@master"
)

func TestGitLabDownloadContent(t *testing.T) {
	r, err := recorder.New("fixtures/content-provider-gitlab")
	if err != nil {
		log.Fatalf("error creating recorder: %v", err)
	}
	defer func(r *recorder.Recorder) { _ = r.Stop() }(r)

	// fake httpClient
	httpClient := r.GetDefaultClient()

	cp, err := NewGitLab(GitLabTestRepository, WithHTTPClient(httpClient))
	assert.NoErrorf(t, err, "error creating GitLab client: %v", err)

	err = cp.DownloadContent()
	assert.NoErrorf(t, err, "error downloading content: %s", err)

	err = cp.Cleanup()
	assert.NoErrorf(t, err, "error cleaning up: %s", err)
}
