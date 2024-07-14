package contentprovider

import (
	"errors"
	"fmt"
	"github.com/huandu/xstrings"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (cp *GitHub) validateRemoteURI(string) error {
	url, ref := parseRemoteURI(cp.remoteURI)
	if !strings.HasPrefix(url, "https://github.com/") {
		return errors.New("invalid github url. The url must start with https://github.com/")
	}
	if ref == "" {
		return errors.New("invalid github url. The url must contain a reference. Example https://github.com/<owner>/<repo>@<ref>")
	}
	cp.repositoryURL = url
	cp.repositoryRef = ref
	return nil
}

func (cp *GitHub) remoteTagExists(tag, authToken string) (bool, error) {
	baseURL := cp.repositoryURIToAPI()
	url := fmt.Sprintf("%s/git/ref/tags/%s", baseURL, tag)
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return false, err
	}
	if authToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", authToken))
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

func (cp *GitHub) repositoryURIToAPI() string {
	_, _, apiPart := xstrings.Partition(cp.repositoryURL, GitHubPrefix)
	return fmt.Sprintf("https://api.github.com/repos/%s", apiPart)
}

func (cp *GitHub) removeGitFolder() error {
	return os.RemoveAll(filepath.Join(cp.workingDir, ".git"))
}
