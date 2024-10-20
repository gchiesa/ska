package localconfigservice

import (
	"fmt"
	"github.com/huandu/xstrings"
	"path/filepath"
	"strings"
	"time"
)

func timeNowUTC() string {
	utcTime := time.Now().UTC()
	timeFormat := "2006-01-02 15:04:05 -0700 MST"
	return utcTime.Format(timeFormat)
}

func makeConfigPath(dirPath string) string {
	return filepath.Join(dirPath, localConfigDirName)
}

func makeConfigFileName(namedConfig string) string {
	if namedConfig == "" {
		namedConfig = localConfigFileNameDefault
	}
	return fmt.Sprintf("%s.%s", namedConfig, localConfigFileNameExt)
}

func hasMultipleConfigurations(dirPath string) bool {
	entries, _ := configEntries(dirPath)
	return len(entries) > 1
}

// configEntries return the list of namedConfigs without any extension
func configEntries(dirPath string) ([]string, error) {
	entries, err := filepath.Glob(dirPath + "/*.yaml")
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, entry := range entries {
		_, _, filename := xstrings.LastPartition(entry, "/")
		configEntry := strings.TrimSuffix(filename, "."+localConfigFileNameExt)
		result = append(result, configEntry)
	}
	return result, nil
}
