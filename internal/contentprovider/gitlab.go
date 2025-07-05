package contentprovider

import (
	"github.com/apex/log"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"net/http"
	"os"
	"path/filepath"
)

type GitLab struct {
	remoteURI          string
	repositoryURL      string
	repositoryRef      string
	repositoryFilePath string
	projectPath        string
	workingDir         string
	gitlabOptions      []gitlab.ClientOptionFunc
	log                *log.Entry
}

type GitLabOption func(*GitLab)

const GitLabPrefix = "https://gitlab.com/"

func NewGitLab(remoteURI string, opt ...GitLabOption) (*GitLab, error) {
	logCtx := log.WithFields(log.Fields{
		"pkg":  "contentprovider",
		"type": "gitlab",
	})
	tmpDir, err := os.MkdirTemp("", workingDirPrefix)
	if err != nil {
		return nil, err
	}

	gl := &GitLab{remoteURI: remoteURI, workingDir: tmpDir, log: logCtx}

	// set options
	for _, opt := range opt {
		opt(gl)
	}
	return gl, nil
}

func WithHTTPClient(httpClient *http.Client) GitLabOption {
	return func(cp *GitLab) {
		cp.gitlabOptions = append(cp.gitlabOptions, gitlab.WithHTTPClient(httpClient))
	}
}

func (cp *GitLab) WorkingDir() string {
	// if filepath is set, we set that one as working directory
	if cp.repositoryFilePath != "" {
		wd := filepath.Join(cp.workingDir, cp.repositoryFilePath)
		cp.log.WithFields(log.Fields{"repositoryFilePath": cp.repositoryFilePath}).Debugf("using repository file path: %s", wd)
		return wd
	}
	return cp.workingDir
}

func (cp *GitLab) Cleanup() error {
	cp.log.WithFields(log.Fields{"workingDir": cp.workingDir}).Debug("removing working dir.")
	err := os.RemoveAll(cp.workingDir)
	return err
}

func (cp *GitLab) RemoteURI() string {
	return cp.remoteURI
}

func (cp *GitLab) DownloadContent() error {
	if err := cp.validateRemoteURI(); err != nil {
		return err
	}

	tmpArchivePath, err := cp.downloadRepoZipArchive()
	if err != nil {
		return err
	}

	// unzip archive
	if err := cp.unzipArchive(cp.workingDir, tmpArchivePath); err != nil {
		return err
	}

	// remove the temporary archive
	if err := os.RemoveAll(tmpArchivePath); err != nil {
		return err
	}

	return nil
}
