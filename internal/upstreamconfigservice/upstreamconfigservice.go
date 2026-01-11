package upstreamconfigservice

import (
	"path/filepath"

	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/utils"
	"gopkg.in/yaml.v2"
)

const upstreamConfigFileName = ".ska-upstream.yaml"

type UpstreamConfigService struct {
	config *config
	log    *log.Entry
}

type UpstreamConfigInput struct {
	Placeholder string   `yaml:"placeholder"`
	Label       string   `yaml:"label"`
	Type        string   `yaml:"type,omitempty"` // "text" (default) or "list"
	Regexp      string   `yaml:"regexp,omitempty"`
	MinLength   int      `yaml:"minLength,omitempty"`
	MaxLength   int      `yaml:"maxLength,omitempty"`
	Help        string   `yaml:"help,omitempty"`
	Default     string   `yaml:"default,omitempty"`
	WriteOnce   bool     `yaml:"writeOnce,omitempty"`
	Items       []string `yaml:"items,omitempty"` // Static list items for list type
}

type SkaConfig struct {
	IgnorePaths []string `yaml:"ignorePaths"`
}

type config struct {
	IgnorePaths []string              `yaml:"ignorePaths"`
	Inputs      []UpstreamConfigInput `yaml:"inputs,omitempty"`
	SkaConfig   SkaConfig             `yaml:"skaConfig,omitempty"`
}

func NewUpstreamConfigService() *UpstreamConfigService {
	logCtx := log.WithFields(log.Fields{
		"pkg": "configuration",
	})
	return &UpstreamConfigService{
		config: &config{},
		log:    logCtx,
	}
}

func (ucs *UpstreamConfigService) LoadFromPath(path string) (*UpstreamConfigService, error) {
	cf := utils.NewConfigFromFile(filepath.Join(path, upstreamConfigFileName))
	configData, err := cf.ReadConfig()
	if err != nil {
		return nil, err
	}

	var cfg config
	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		return nil, err
	}

	ucs.config = &cfg
	return ucs, nil
}

func (ucs *UpstreamConfigService) UpstreamIgnorePaths() []string {
	return ucs.config.IgnorePaths
}

func (ucs *UpstreamConfigService) SkaConfigIgnorePaths() []string {
	return ucs.config.SkaConfig.IgnorePaths
}

func (ucs *UpstreamConfigService) GetInputs() []UpstreamConfigInput {
	return ucs.config.Inputs
}
