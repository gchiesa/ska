package contentprovider

import (
	"github.com/huandu/xstrings"
)

func parseRemoteURI(uri string) (url, tag string) {
	url, _, tag = xstrings.Partition(uri, "@")
	return url, tag
}
