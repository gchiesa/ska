package configuration

import (
	"errors"
	"fmt"
	"github.com/huandu/xstrings"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"slices"
	"time"
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
	newPaths := append(cs.app.Config.IgnorePaths, ignorePaths...)
	slices.Sort(newPaths)
	cs.app.Config.IgnorePaths = slices.Compact(newPaths)
	return cs
}

func (cs *LocalConfigService) WithVariables(variables map[string]interface{}) *LocalConfigService {
	cs.app.State.Variables = variables
	return cs
}

func (cs *LocalConfigService) WriteConfig(dirPath string) error {
	if err := os.MkdirAll(makeConfigPath(dirPath), 0o755); err != nil {
		return err
	}

	cf := NewConfigFromFile(filepath.Join(makeConfigPath(dirPath), makeConfigFileName(cs.namedConfig)))

	// get the time utc now in format "2006-01-02 15:04:05.999999999 -0700 MST"
	cs.app.State.LastUpdate = TimeNowUTC()

	configData, err := yaml.Marshal(cs.app)
	if err != nil {
		return err
	}

	if err := cf.WriteConfig(configData); err != nil {
		return err
	}
	return nil
}

func (cs *LocalConfigService) ReadValidConfig(dirPath string) error {
	if hasMultipleConfigurations(makeConfigPath(dirPath)) && cs.namedConfig == "" {
		return ErrNoConfigSpecified
	}

	cf := NewConfigFromFile(filepath.Join(makeConfigPath(dirPath), makeConfigFileName(cs.namedConfig)))
	return cs.LoadConfig(cf)
}

func (cs *LocalConfigService) LoadConfig(cf *ConfigFile) error {
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
	configFileName := makeConfigFileName(cs.namedConfig)
	existingConfigs, _ := configEntries(makeConfigPath(dirPath))
	return slices.Contains(existingConfigs, configFileName)
}

func TimeNowUTC() string {
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

func configEntries(dirPath string) ([]string, error) {
	entries, err := filepath.Glob(dirPath + "/*.yaml")
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, entry := range entries {
		_, _, filename := xstrings.LastPartition(entry, "/")
		result = append(result, filename)
	}
	return result, nil
}
