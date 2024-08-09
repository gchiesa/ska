package contentprovider

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// const GitLabTestPublicRepo = "https://gitlab.com/gchiesa/s3vaultlib@master"
const GitLabTestPublicRepo = "https://gitlab.com/gchiesa/test@master"

func TestGitLab_DownloadPublicRepo(t *testing.T) {
	var err error
	cp, err := NewGitLab(GitLabTestPublicRepo)
	assert.NoErrorf(t, err, "error creating gitlab provider: %s", err)

	err = cp.DownloadContent()
	assert.NoErrorf(t, err, "error downloading content: %s", err)

	err = cp.Cleanup()
	assert.NoErrorf(t, err, "error cleaning up: %s", err)
}
