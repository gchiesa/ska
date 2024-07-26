package tui

import (
	"fmt"
	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/configuration"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"regexp"
)

const skaInteractiveFileName = ".ska-interactive.yaml"

type SkaInteractiveService struct {
	templateURI string
	formTitle   string
	formConfig  *InteractiveForm
	variables   map[string]string
	log         *log.Entry
}

type InteractiveInput struct {
	Placeholder string `yaml:"placeholder"`
	Label       string `yaml:"label"`
	RegExp      string `yaml:"regexp,omitempty"`
	MinLength   int    `yaml:"minLength,omitempty"`
	MaxLength   int    `yaml:"maxLength,omitempty"`
	Default     string `yaml:"default,omitempty"`
	Help        string `yaml:"help,omitempty"`
	Value       string
}
type InteractiveForm struct {
	Inputs []InteractiveInput `yaml:"inputs"`
}

type Variables map[string]string

func NewSkaInteractiveService(templateURI, formTitle string) *SkaInteractiveService {
	return &SkaInteractiveService{
		templateURI: templateURI,
		formTitle:   formTitle,
		log:         log.WithFields(log.Fields{"pkg": "skaffolder"}),
	}
}

func (s *SkaInteractiveService) ShouldRun() bool {
	interactiveConfigFilePath := filepath.Join(s.templateURI, skaInteractiveFileName)

	// check if file exists
	if !fileExists(interactiveConfigFilePath) {
		s.log.WithFields(log.Fields{"templateURI": s.templateURI, "interactiveConfig": skaInteractiveFileName}).Debug("no interactive config found.")
		return false
	}

	// check if we can load it
	cf := configuration.NewConfigFromFile(interactiveConfigFilePath)

	configData, err := cf.ReadConfig()
	if err != nil {
		s.log.WithError(err).WithFields(log.Fields{"interactiveConfigFilePath": interactiveConfigFilePath}).Error("could not read interactive config.")
		return false
	}

	var cfg InteractiveForm
	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		s.log.WithError(err).WithFields(log.Fields{"interactiveConfigFilePath": interactiveConfigFilePath}).Error("could not unmarshal interactive config.")
		return false
	}

	s.formConfig = &cfg

	if len(s.formConfig.Inputs) == 0 {
		s.log.Info("no inputs in the interactive config.")
		return false
	}
	return true
}

func (s *SkaInteractiveService) Run() error {
	if !s.ShouldRun() {
		s.log.Info("skipping interactive config")
		return nil
	}
	s.disableWithLoggingInvalidRegExp()

	tui := NewModelFromInteractiveForm(*s.formConfig, s.formTitle)

	if err := tui.Execute(); err != nil {
		return err
	}
	if tui.exitWithCtrlC {
		return fmt.Errorf("cancelled by user")
	}
	s.variables = tui.GetVariablesForInteractiveForm(*s.formConfig)
	return nil
}

func (s *SkaInteractiveService) Variables() map[string]string {
	return s.variables
}

func (s *SkaInteractiveService) SetDefaults(variables map[string]string) {
	for i := range s.formConfig.Inputs {
		if v, ok := variables[s.formConfig.Inputs[i].Placeholder]; ok {
			s.formConfig.Inputs[i].Default = v
		}
	}
}

func (s *SkaInteractiveService) disableWithLoggingInvalidRegExp() {
	for i := range s.formConfig.Inputs {
		if _, err := regexp.Compile(s.formConfig.Inputs[i].RegExp); err != nil {
			s.log.WithFields(log.Fields{"validation": s.formConfig.Inputs[i].RegExp}).Warnf("RegExp expression is invalid. Error: %s. Ignoring validation.", err)
			s.formConfig.Inputs[i].RegExp = ""
		}
	}
}
