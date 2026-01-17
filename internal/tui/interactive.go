package tui

import (
	"fmt"
	"regexp"

	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/upstreamconfigservice"
)

type SkaInteractiveService struct {
	formTitle        string
	formConfig       *InteractiveForm
	formShowBanner   bool
	variables        map[string]string
	writeOnceEnabled bool // writeOnceEnabled enable the readonly management
	log              *log.Entry
}

// InputType defines the type of input (text or list)
type InputType string

const (
	InputTypeText InputType = "text"
	InputTypeList InputType = "list"
)

// ItemsFunc is a function that returns items for a list input
type ItemsFunc func() []string

type InteractiveInput struct {
	Placeholder string    `yaml:"placeholder"`
	Label       string    `yaml:"label"`
	Type        InputType `yaml:"type,omitempty"` // "text" (default) or "list"
	RegExp      string    `yaml:"regexp,omitempty"`
	MinLength   int       `yaml:"minLength,omitempty"`
	MaxLength   int       `yaml:"maxLength,omitempty"`
	Default     string    `yaml:"default,omitempty"`
	WriteOnce   bool      `yaml:"writeOnce,omitempty"`
	Help        string    `yaml:"help,omitempty"`
	Value       string
	Items       []string  `yaml:"items,omitempty"` // Static list items for list type
	ItemsFunc   ItemsFunc `yaml:"-"`               // Function to generate list items dynamically (not serialized)
}
type InteractiveForm struct {
	Inputs []InteractiveInput `yaml:"inputs"`
}

type Variables map[string]string

func NewSkaInteractiveService(formTitle string, inputs []upstreamconfigservice.UpstreamConfigInput) *SkaInteractiveService {
	var interactiveInputs []InteractiveInput

	for _, i := range inputs {
		inputType := InputType(i.Type)
		if inputType == "" {
			inputType = InputTypeText
		}

		input := &InteractiveInput{
			Placeholder: i.Placeholder,
			Label:       i.Label,
			Type:        inputType,
			RegExp:      i.Regexp,
			MinLength:   i.MinLength,
			MaxLength:   i.MaxLength,
			Default:     i.Default,
			WriteOnce:   i.WriteOnce,
			Help:        i.Help,
			Items:       i.Items,
		}
		interactiveInputs = append(interactiveInputs, *input)
	}

	return &SkaInteractiveService{
		formTitle:  formTitle,
		formConfig: &InteractiveForm{Inputs: interactiveInputs},
		log:        log.WithFields(log.Fields{"pkg": "skaffolder"}),
	}
}

func (s *SkaInteractiveService) SetWriteOnce(isEnabled bool) *SkaInteractiveService {
	s.writeOnceEnabled = isEnabled
	return s
}

func (s *SkaInteractiveService) ShouldRun() bool {
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

	tui := NewModelFromInteractiveForm(*s.formConfig, s.formTitle, s.formShowBanner).
		SetWriteOnce(s.writeOnceEnabled)

	if err := tui.Banner(); err != nil {
		return err
	}
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

func (s *SkaInteractiveService) SetShowBanner(enabled bool) *SkaInteractiveService {
	s.formShowBanner = enabled
	return s
}

func (s *SkaInteractiveService) disableWithLoggingInvalidRegExp() {
	for i := range s.formConfig.Inputs {
		if _, err := regexp.Compile(s.formConfig.Inputs[i].RegExp); err != nil {
			s.log.WithFields(log.Fields{"validation": s.formConfig.Inputs[i].RegExp}).Warnf("the RegExp expression is invalid. Error: %s. Ignoring validation.", err)
			s.formConfig.Inputs[i].RegExp = ""
		}
	}
}
