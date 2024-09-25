package configuration

import (
	"github.com/apex/log"
	"os"
)

type ConfigFile struct {
	filePath string
	log      *log.Entry
}

func NewConfigFromFile(filePath string) *ConfigFile {
	logCtx := log.WithFields(log.Fields{
		"pkg": "configuration",
	})
	return &ConfigFile{
		filePath: filePath,
		log:      logCtx,
	}
}

func (cf *ConfigFile) GetFilePath() string {
	return cf.filePath
}

func (cf *ConfigFile) WriteConfig(configData []byte) error {
	cf.log.WithFields(log.Fields{"filePath": cf.filePath}).Debug("writing configuration to file")
	if err := os.WriteFile(cf.filePath, configData, 0o644); err != nil {
		return err
	}
	return nil
}

func (cf *ConfigFile) ReadConfig() ([]byte, error) {
	cf.log.WithFields(log.Fields{"filePath": cf.filePath}).Debug("reading configuration from file")
	return os.ReadFile(cf.filePath)
}
