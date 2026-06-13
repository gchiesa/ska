package contentprovider

import (
	"fmt"
	"log/slog"
	"strings"
)

// ContentProviderOption configures a content provider returned by ByURI.
type ContentProviderOption func(*contentProviderConfig)

type contentProviderConfig struct {
	logger *slog.Logger
}

// WithLogger injects a *slog.Logger into the content provider.
// Each provider implementation will add its own "pkg" / "type" fields.
func WithLogger(logger *slog.Logger) ContentProviderOption {
	return func(c *contentProviderConfig) {
		c.logger = logger
	}
}

func ByURI(uri string, opts ...ContentProviderOption) (RemoteContentProvider, error) {
	cfg := &contentProviderConfig{logger: slog.Default()}
	for _, opt := range opts {
		opt(cfg)
	}

	var err error
	var contentProvider RemoteContentProvider
	switch {
	case strings.HasPrefix(uri, GitHubPrefix):
		contentProvider, err = NewGitHub(uri, cfg.logger)
	case strings.HasPrefix(uri, GitLabPrefix):
		contentProvider, err = NewGitLab(uri, cfg.logger)
	case strings.HasPrefix(uri, LocalPathPrefix):
		contentProvider, err = NewLocalPath(strings.TrimPrefix(uri, LocalPathPrefix), cfg.logger)
	default:
		return nil, fmt.Errorf("unsupported uri: %s", uri)
	}

	if err != nil {
		return nil, err
	}
	return contentProvider, nil
}
