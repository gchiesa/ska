package configuration

import (
	"gopkg.in/yaml.v2"
	"time"
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
	app *appCfg
}

func NewLocalConfigService() *LocalConfigService {
	configBlock := &ConfigBlock{}
	stateBlock := &StateBlock{}
	appConfiguration := &appCfg{
		Config: *configBlock,
		State:  *stateBlock,
	}
	return &LocalConfigService{app: appConfiguration}
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

func (cs *LocalConfigService) GetIgnorePaths() []string {
	return cs.app.Config.IgnorePaths
}

func (cs *LocalConfigService) WithVariables(variables map[string]interface{}) *LocalConfigService {
	cs.app.State.Variables = variables
	return cs
}

func (cs *LocalConfigService) WriteConfig(dirPath string) error {
	cf := NewConfigFromDirectory(dirPath)

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

func (cs *LocalConfigService) ReadConfig(dirPath string) error {
	cf := NewConfigFromDirectory(dirPath)

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

func TimeNowUTC() string {
	utcTime := time.Now().UTC()
	timeFormat := "2006-01-02 15:04:05 -0700 MST"
	return utcTime.Format(timeFormat)
}
