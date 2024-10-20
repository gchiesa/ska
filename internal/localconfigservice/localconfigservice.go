package localconfigservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gchiesa/ska/internal/utils"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"slices"
)

var ErrNoConfigSpecified = errors.New("no configuration specified and multiple configurations present")

const (
	localConfigDirName         = ".ska-config"
	localConfigFileNameDefault = "default"
	localConfigFileNameExt     = "yaml"
)

type ConfigBlock struct {
	BlueprintURI string   `yaml:"blueprintURI"`
	IgnorePaths  []string `yaml:"ignorePaths"`
}

type StateBlock struct {
	LastUpdate string                 `yaml:"lastUpdate"`
	Variables  map[string]interface{} `yaml:"variables"`
}

type appCfg struct {
	Config ConfigBlock `yaml:"config"`
	State  StateBlock  `yaml:"state"`
}

type LocalConfigService struct {
	namedConfig string
	app         *appCfg
}

func NewLocalConfigService(namedConfig string) *LocalConfigService {
	configBlock := &ConfigBlock{}
	stateBlock := &StateBlock{}
	appConfiguration := &appCfg{
		Config: *configBlock,
		State:  *stateBlock,
	}
	return &LocalConfigService{namedConfig: namedConfig, app: appConfiguration}
}

func (cs *LocalConfigService) NamedConfig() string {
	if cs.namedConfig != "" {
		return cs.namedConfig
	}
	return localConfigFileNameDefault
}

func (cs *LocalConfigService) BlueprintUpstream() string {
	return cs.app.Config.BlueprintURI
}

func (cs *LocalConfigService) Variables() map[string]interface{} {
	return cs.app.State.Variables
}

func (cs *LocalConfigService) WithBlueprintUpstream(bpURI string) *LocalConfigService {
	cs.app.Config.BlueprintURI = bpURI
	return cs
}

func (cs *LocalConfigService) ProcessAllFiles() bool {
	return len(cs.app.Config.IgnorePaths) == 0
}

func (cs *LocalConfigService) WithExcludeMatchingFiles(ignorePaths []string) *LocalConfigService {
	cs.app.Config.IgnorePaths = ignorePaths
	return cs
}

func (cs *LocalConfigService) IgnorePaths() []string {
	return cs.app.Config.IgnorePaths
}

func (cs *LocalConfigService) WithIgnorePaths(ignorePaths []string) *LocalConfigService {
	cs.app.Config.IgnorePaths = ignorePaths
	return cs
}

func (cs *LocalConfigService) WithExtendIgnorePaths(ignorePaths []string) *LocalConfigService {
	newPaths := append(cs.app.Config.IgnorePaths, ignorePaths...) //nolint:gocritic
	slices.Sort(newPaths)
	cs.app.Config.IgnorePaths = slices.Compact(newPaths)
	return cs
}

func (cs *LocalConfigService) WithVariables(variables map[string]interface{}) *LocalConfigService {
	cs.app.State.Variables = variables
	return cs
}

func (cs *LocalConfigService) RenameNamedConfig(dirPath, namedConfig string) error {
	if cs.ConfigExistsWithNamedConfig(dirPath, namedConfig) {
		return fmt.Errorf("cannot rename %s to %s because named configuration already exists on path: %s",
			cs.namedConfig, namedConfig, dirPath)
	}

	cs.namedConfig = namedConfig

	if err := cs.WriteConfig(dirPath); err != nil {
		return err
	}
	return cs.DeleteConfig(dirPath)
}

func (cs *LocalConfigService) WriteConfig(dirPath string) error {
	if err := os.MkdirAll(makeConfigPath(dirPath), 0o755); err != nil {
		return err
	}

	cf := utils.NewConfigFromFile(filepath.Join(makeConfigPath(dirPath), makeConfigFileName(cs.namedConfig)))

	// get the time utc now in format "2006-01-02 15:04:05.999999999 -0700 MST"
	cs.app.State.LastUpdate = timeNowUTC()

	configData, err := yaml.Marshal(cs.app)
	if err != nil {
		return err
	}

	if err := cf.WriteConfig(configData); err != nil {
		return err
	}
	return nil
}

func (cs *LocalConfigService) ToJSON() ([]byte, error) {
	configData, err := json.Marshal(cs.app)
	if err != nil {
		return nil, err
	}
	return configData, nil
}

func (cs *LocalConfigService) DeleteConfig(dirPath string) error {
	return os.RemoveAll(makeConfigPath(dirPath))
}

func (cs *LocalConfigService) ReadValidConfig(dirPath string) error {
	if hasMultipleConfigurations(makeConfigPath(dirPath)) && cs.namedConfig == "" {
		return ErrNoConfigSpecified
	}

	cf := utils.NewConfigFromFile(filepath.Join(makeConfigPath(dirPath), makeConfigFileName(cs.namedConfig)))
	return cs.LoadConfig(cf)
}

func (cs *LocalConfigService) LoadConfig(cf *utils.ConfigFile) error {
	configData, err := cf.ReadConfig()
	if err != nil {
		return err
	}

	var cfg appCfg
	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		return err
	}
	cs.app = &cfg
	return nil
}

func (cs *LocalConfigService) ConfigExists(dirPath string) bool {
	return cs.ConfigExistsWithNamedConfig(dirPath, cs.namedConfig)
}

func (cs *LocalConfigService) ConfigExistsWithNamedConfig(dirPath, namedConfig string) bool {
	existingConfigs, _ := configEntries(makeConfigPath(dirPath))
	return slices.Contains(existingConfigs, namedConfig)
}
