package contentprovider

import (
	"archive/zip"
	"fmt"
	"github.com/apex/log"
	"github.com/xanzy/go-gitlab"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type GitLab struct {
	remoteURI     string
	repositoryURL string
	repositoryRef string
	projectPath   string
	workingDir    string
	log           *log.Entry
}

const GitLabPrefix = "https://gitlab.com/"

func NewGitLab(remoteURI string) (*GitLab, error) {
	logCtx := log.WithFields(log.Fields{
		"pkg":  "contentprovider",
		"type": "gitlab",
	})
	tmpDir, err := os.MkdirTemp("", workingDirPrefix)
	if err != nil {
		return nil, err
	}
	return &GitLab{remoteURI: remoteURI, workingDir: tmpDir, log: logCtx}, nil
}

func (cp *GitLab) WorkingDir() string {
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

func (cp *GitLab) downloadRepoZipArchive() (zipArchive string, err error) {
	token := os.Getenv("GITLAB_PRIVATE_TOKEN")

	gitlabClient, err := gitlab.NewClient(token)
	if err != nil {
		return "", err
	}

	tmpArchive, err := os.CreateTemp(os.TempDir(), "gitlab-repo-")
	defer func(f *os.File) { _ = f.Close() }(tmpArchive)

	if err != nil {
		return "", err
	}

	var archiveFormat = "zip"
	archiveOptions := &gitlab.ArchiveOptions{
		Format: &archiveFormat,
		SHA:    &cp.repositoryRef,
	}
	data, resp, err := gitlabClient.Repositories.Archive(cp.projectPath, archiveOptions, gitlab.WithToken(gitlab.PrivateToken, token))

	if resp.StatusCode == 404 {
		return "", fmt.Errorf("repository not found, perhaps you need to setup your private token for GitLab by exporting the environment variable GITLAB_PRIVATE_TOKEN. Status Code: %d", resp.StatusCode)
	}

	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("invalid status code returned from GitLab API. Status Code: %d", resp.StatusCode)
	}

	// write data to tmp archive
	if _, err := tmpArchive.Write(data); err != nil {
		return "", err
	}

	return tmpArchive.Name(), nil
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

func (cp *GitLab) unzipArchive(dst, archivePath string) error {
	archive, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer func(*zip.ReadCloser) { _ = archive.Close() }(archive)

	rootPath := ""
	for _, item := range archive.File {
		dstFilePath := filepath.Join(dst, item.Name) //nolint:gosec

		if !strings.HasPrefix(dstFilePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", dstFilePath)
		}
		if item.FileInfo().IsDir() {
			if err := os.MkdirAll(dstFilePath, os.ModePerm); err != nil {
				return err
			}
			if rootPath == "" {
				rootPath = item.Name
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dstFilePath), os.ModePerm); err != nil {
			return err
		}

		// open the dstFilePath as file
		dstFile, err := os.OpenFile(dstFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, item.Mode())
		if err != nil {
			panic(err)
		}

		compressedFile, err := item.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, compressedFile); err != nil { //nolint:gosec
			return err
		}
		_ = dstFile.Close()
		_ = compressedFile.Close()
	}
	cp.workingDir = filepath.Join(dst, rootPath)
	return nil
}
