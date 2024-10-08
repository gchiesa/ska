package contentprovider

import (
	"fmt"
	"strings"
)

func ByURI(uri string) (RemoteContentProvider, error) {
	var err error
	var contentProvider RemoteContentProvider
	switch {
	case strings.HasPrefix(uri, GitHubPrefix):
		contentProvider, err = NewGitHub(uri)
	case strings.HasPrefix(uri, GitLabPrefix):
		contentProvider, err = NewGitLab(uri)
	case strings.HasPrefix(uri, LocalPathPrefix):
		contentProvider, err = NewLocalPath(strings.TrimPrefix(uri, LocalPathPrefix))
	default:
		return nil, fmt.Errorf("unsupported uri: %s", uri)
	}

	if err != nil {
		return nil, err
	}
	return contentProvider, nil
}
