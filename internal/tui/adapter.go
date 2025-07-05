package tui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
)

func NewModelFromInteractiveForm(iForm InteractiveForm, header string, showBanner bool) *Model {
	m := &Model{
		header:     header,
		showBanner: showBanner,
		inputs:     make([]textinput.Model, len(iForm.Inputs)),
	}

	for i := range iForm.Inputs {
		t := textinput.New()

		// Prompt
		t.Prompt = fmt.Sprintf("%s: ", iForm.Inputs[i].Label)
		// Placeholder
		t.Placeholder = iForm.Inputs[i].Help
		t.PlaceholderStyle = helpStyle
		// Validation
		t.Validate = validator(iForm.Inputs[i].WriteOnce, iForm.Inputs[i].MinLength, iForm.Inputs[i].RegExp, iForm.Inputs[i].Default)
		if iForm.Inputs[i].MaxLength > 0 {
			t.CharLimit = iForm.Inputs[i].MaxLength
		}
		// Default
		if iForm.Inputs[i].Default != "" {
			t.SetValue(iForm.Inputs[i].Default)
		}
		// First Item
		if i == 0 {
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = noStyle
		}
		m.inputs[i] = t
	}
	return m
}

func validator(writeOnce bool, minLength int, regexpString, oldValue string) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" && minLength > 0 {
			return fmt.Errorf("value cannot be empty")
		}
		if len(s) < minLength {
			return fmt.Errorf("value is too short (min length: %d)", minLength)
		}
		if strings.TrimSpace(regexpString) != "" {
			re := regexp.MustCompile(regexpString)
			if !re.MatchString(s) {
				return fmt.Errorf("invalid value. It should match %s", regexpString)
			}
		}
		if writeOnce && oldValue != "" {
			if s != oldValue {
				return fmt.Errorf("value cannot be changed, please change it back to '%s'", oldValue)
			}
		}
		return nil
	}
}
