package contentprovider

import (
	"fmt"
	"github.com/huandu/xstrings"
	"strings"
)

const (
	logPkg             = "contentprovider"
	logFieldPkg        = "pkg"
	logFieldType       = "type"
	logFieldWorkingDir = "workingDir"
)

func parseRemoteURIV2(uri string) (url, filePath, tag string) {
	urlWithPath, _, tag := xstrings.Partition(uri, "@")
	urlWithNoSchema := strings.TrimPrefix(urlWithPath, "https://")
	urlWithNoPath, _, path := xstrings.Partition(urlWithNoSchema, "//")
	return fmt.Sprintf("https://%s", urlWithNoPath), path, tag
}
