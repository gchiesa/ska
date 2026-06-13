package utils

import (
	"fmt"
	"log/slog"
	"os"
)

type ConfigFile struct {
	filePath string
	log      *slog.Logger
}

func NewConfigFromFile(filePath string, logger *slog.Logger) *ConfigFile {
	if logger == nil {
		logger = slog.Default()
	}
	return &ConfigFile{
		filePath: filePath,
		log:      logger.With("pkg", "configuration"),
	}
}

func (cf *ConfigFile) GetFilePath() string {
	return cf.filePath
}

func (cf *ConfigFile) WriteConfig(configData []byte) error {
	cf.log.With("filePath", cf.filePath).Debug("writing configuration to file")
	if err := os.WriteFile(cf.filePath, configData, 0o644); err != nil {
		return fmt.Errorf("writing configuration: %w", err)
	}
	return nil
}

func (cf *ConfigFile) ReadConfig() ([]byte, error) {
	cf.log.With("filePath", cf.filePath).Debug("reading configuration from file")
	return os.ReadFile(cf.filePath)
}
