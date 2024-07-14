package configuration

import (
	"gopkg.in/yaml.v2"
)

type ConfigBlock struct {
	BlueprintUpstream    string   `yaml:"blueprintUpstream"`
	ExcludeMatchingFiles []string `yaml:"excludeMatchingFiles"`
	IncludeMatchingFiles []string `yaml:"includeMatchingFiles"`
}

type StateBlock struct {
	Variables map[string]interface{}
}

type appCfg struct {
	Config ConfigBlock `yaml:"config"`
	State  StateBlock  `yaml:"state"`
}

type ConfigService struct {
	app *appCfg
}

func NewConfigService() *ConfigService {
	configBlock := &ConfigBlock{}
	stateBlock := &StateBlock{}
	appConfiguration := &appCfg{
		Config: *configBlock,
		State:  *stateBlock,
	}
	return &ConfigService{app: appConfiguration}
}

func (cs *ConfigService) WithBlueprintUpstream(bpURI string) *ConfigService {
	cs.app.Config.BlueprintUpstream = bpURI
	return cs
}

func (cs *ConfigService) ProcessAllFiles() bool {
	return len(cs.app.Config.ExcludeMatchingFiles) == 0 && len(cs.app.Config.IncludeMatchingFiles) == 0
}

func (cs *ConfigService) WithExcludeMatchingFiles(excludeMatchingFiles []string) *ConfigService {
	cs.app.Config.ExcludeMatchingFiles = excludeMatchingFiles
	return cs
}

func (cs *ConfigService) WithIncludeMatchingFiles(includeMatchingFiles []string) *ConfigService {
	cs.app.Config.IncludeMatchingFiles = includeMatchingFiles
	return cs
}

func (cs *ConfigService) WithVariables(variables map[string]interface{}) *ConfigService {
	cs.app.State.Variables = variables
	return cs
}

func (cs *ConfigService) WriteConfig(dirPath string) error {
	cf := NewConfigFromDirectory(dirPath)

	configData, err := yaml.Marshal(cs.app)
	if err != nil {
		return err
	}

	if err := cf.WriteConfig(configData); err != nil {
		return err
	}
	return nil
}

func (cs *ConfigService) ReadConfig(dirPath string) ([]byte, error) {
	cf := NewConfigFromDirectory(dirPath)

	configData, err := cf.ReadConfig()
	if err != nil {
		return nil, err
	}

	var cfg appCfg
	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		return nil, err
	}
	cs.app = &cfg
	return cf.ReadConfig()
}
