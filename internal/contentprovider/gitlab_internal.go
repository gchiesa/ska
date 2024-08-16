package contentprovider

import (
	"fmt"
	"strings"
)

func (cp *GitLab) validateRemoteURI() error {
	url, ref := parseRemoteURI(cp.remoteURI)
	if !strings.HasPrefix(url, GitLabPrefix) {
		return fmt.Errorf("invalid github url. The url must start with: %s", GitLabPrefix)
	}
	if ref == "" {
		return fmt.Errorf("invalid github url. The url must contain a reference. Example %s/<namespace>/<repo>@<ref>", GitLabPrefix)
	}
	cp.repositoryURL = url
	cp.repositoryRef = ref
	cp.projectPath = strings.TrimPrefix(cp.repositoryURL, GitLabPrefix)
	return nil
}
