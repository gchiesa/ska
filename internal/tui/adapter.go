package tui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	header           string
	showBanner       bool
	focusIndex       int
	entries          []inputEntry    // new unified entry system
	readonlyLabelMap map[string]bool // readonlyLabelMap map of labels which are readonly
	err              error
	exitWithCtrlC    bool
}

// inputEntry represents a single input in the form (either text or list)
type inputEntry struct {
	inputType InputType
	textInput textinput.Model
	listModel list.Model
	label     string
	prompt    string // formatted prompt (e.g., "Label      : ")
	selected  string // selected value for list inputs
}

func NewModelFromInteractiveForm(iForm InteractiveForm, header string, showBanner bool) *Model {
	m := &Model{
		header:           header,
		showBanner:       showBanner,
		entries:          make([]inputEntry, len(iForm.Inputs)),
		readonlyLabelMap: make(map[string]bool, len(iForm.Inputs)),
	}

	// calculate max placeholder length
	var maxPromptLength int
	for i := range iForm.Inputs {
		maxPromptLength = max(maxPromptLength, len(iForm.Inputs[i].Label))
	}

	promptFormat := fmt.Sprintf("%%-%ds: ", maxPromptLength)
	for i := range iForm.Inputs {
		inputType := iForm.Inputs[i].Type
		if inputType == "" {
			inputType = InputTypeText
		}

		entry := inputEntry{
			inputType: inputType,
			label:     iForm.Inputs[i].Label,
			prompt:    fmt.Sprintf(promptFormat, iForm.Inputs[i].Label),
		}

		// fill the map with readonly labels
		if iForm.Inputs[i].WriteOnce {
			m.readonlyLabelMap[entry.label] = false
		}

		if inputType == InputTypeList {
			// Create list input
			entry.listModel = createListWidget(iForm.Inputs[i])
			// Set default selection if provided
			if iForm.Inputs[i].Default != "" {
				entry.selected = iForm.Inputs[i].Default
			}
		} else {
			// Create text input (default)
			t := textinput.New()

			// Prompt
			t.Prompt = fmt.Sprintf(promptFormat, iForm.Inputs[i].Label)

			// Placeholder
			t.Placeholder = iForm.Inputs[i].Help
			t.PlaceholderStyle = helpStyle
			// Validation
			t.Validate = validator(iForm.Inputs[i].MinLength, iForm.Inputs[i].RegExp)
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
			entry.textInput = t
		}
		m.entries[i] = entry
	}
	return m
}

func (m *Model) SetWriteOnce(isEnabled bool) *Model {
	// update the state on the labels
	for k := range m.readonlyLabelMap {
		m.readonlyLabelMap[k] = isEnabled
	}
	return m
}

// Banner shows the SKA banner either with text or graphical interface
func (m *Model) Banner() error {
	if !m.showBanner {
		return nil
	}

	if !CanUseGraphic() {
		println(TextBanner())
		return nil
	}
	return GraphicalBanner()
}

func validator(minLength int, regexpString string) func(string) error {
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
		return nil
	}
}
