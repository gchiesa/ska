package content_provider

import (
	"github.com/huandu/xstrings"
)

func parseRemoteURI(uri string) (string, string) {
	url, _, tag := xstrings.Partition(uri, "@")
	return url, tag
}
