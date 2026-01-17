package tui

import (
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg error
)

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#416767"}
	highlight = lipgloss.AdaptiveColor{Light: "#83ADF4", Dark: "#83ADF4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	good      = lipgloss.AdaptiveColor{Light: "#32a71d", Dark: "#32a71d"}
	bad       = lipgloss.AdaptiveColor{Light: "#CE1E00", Dark: "#CE1E00"}

	headerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			PaddingLeft(2).PaddingRight(2).Foreground(special)

	focusedStyle = lipgloss.NewStyle().Bold(true).Foreground(highlight)
	blurredStyle = lipgloss.NewStyle().Bold(false).Foreground(subtle)
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle.AlignHorizontal(lipgloss.Right).Italic(true)
	errorStyle   = lipgloss.NewStyle().Foreground(bad).MarginTop(2).MarginBottom(1)
	goodTick     = lipgloss.NewStyle().Foreground(good)
	badTick      = lipgloss.NewStyle().Foreground(bad)

	listLabelStyle = lipgloss.NewStyle().Bold(true).Foreground(highlight)

	// listBoxStyle creates a rounded border around the list dropdown
	listBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(0, 1)
)

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Update all entries based on their type
	for i := range m.entries {
		// we process the update on the list only if the list is opened
		if i == m.focusIndex && m.entries[i].inputType == InputTypeList {
			var cmd tea.Cmd
			m.entries[i].listModel, cmd = m.entries[i].listModel.Update(msg)
			cmds = append(cmds, cmd)
			// Update selected value if item is selected
			if item := m.entries[i].listModel.SelectedItem(); item != nil {
				m.entries[i].selected = item.(listItem).title
			}
		} else {
			var cmd tea.Cmd
			m.entries[i].textInput, cmd = m.entries[i].textInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		currentEntry := &m.entries[m.focusIndex]

		switch msg.Type {
		case tea.KeyCtrlS:
			// Validate all entries
			var i int
			for i = range m.entries {
				if m.entries[i].inputType == InputTypeText {
					if err := m.entries[i].textInput.Validate(m.entries[i].textInput.Value()); err != nil {
						m.err = err
						break
					}
				}
				// List inputs are always valid (selection is required)
				if m.entries[i].inputType == InputTypeList && m.entries[i].selected == "" {
					m.err = fmt.Errorf("please select an item for %s", m.entries[i].label)
					break
				}
			}
			if m.err != nil {
				m.focusEntry(i)
			} else {
				return m, tea.Quit
			}

		case tea.KeyEnter:
			// For list inputs, Enter selects and moves to next
			if currentEntry.inputType == InputTypeList {
				if item := currentEntry.listModel.SelectedItem(); item != nil {
					currentEntry.selected = item.(listItem).title
				}
				m.nextEntryIfNoError()
			} else {
				// For text inputs, Enter moves to next
				if m.focusIndex == len(m.entries)-1 {
					// Last entry, check if we should quit
					return m, tea.Quit
				}
				m.nextEntryIfNoError()
			}

		case tea.KeyTab:
			m.nextEntryIfNoError()

		case tea.KeyShiftTab:
			m.prevEntryIfNoError()

		case tea.KeyCtrlC, tea.KeyEsc:
			if currentEntry.inputType == InputTypeList {
				// if we are here we just want to exit from the list and go to the previous item
				// we need to delete the cmd created by the list update above that would otherwise exit the application
				cmds[m.focusIndex] = nil
				m.prevEntryIfNoError()
			} else {
				m.exitWithCtrlC = true
				return m, tea.Quit
			}

		case tea.KeyUp, tea.KeyDown:
			// For list inputs, let the list handle up/down
			if currentEntry.inputType == InputTypeList {
				// Already handled above in the list update
			} else {
				// For text inputs, navigate between entries
				if msg.Type == tea.KeyDown {
					m.nextEntryIfNoError()
				} else {
					m.prevEntryIfNoError()
				}
			}
		}

		// Update focus styles
		m.updateFocusStyles()

	case errMsg:
		m.err = msg
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) updateFocusStyles() {
	for i := range m.entries {
		if m.entries[i].inputType == InputTypeText {
			if i == m.focusIndex {
				m.entries[i].textInput.PromptStyle = focusedStyle
				m.entries[i].textInput.Focus()
			} else {
				m.entries[i].textInput.PromptStyle = noStyle
				m.entries[i].textInput.Blur()
			}
		}
	}
}

func (m *Model) View() string {
	var builder strings.Builder

	builder.WriteString(headerStyle.Render(m.header))
	builder.WriteRune('\n')
	builder.WriteString("Please fill the required fields below:\n\n")

	for i := range m.entries {
		entry := &m.entries[i]

		if entry.inputType == InputTypeList {
			// Render list input
			hasSelection := entry.selected != ""
			if hasSelection {
				builder.WriteString(goodTick.Render("✔"))
			} else {
				builder.WriteString(badTick.Render("✖"))
			}

			// if the list is focused
			if i == m.focusIndex {
				// Show the label and list dropdown when focused
				builder.WriteString(fmt.Sprintf(" %s", listLabelStyle.Render(entry.prompt)))

				// Calculate the indentation of the list box to align with where the value would appear
				// "✔ " = 2 chars, then prompt length
				indentRemainingBoxLines := strings.Repeat(" ", 2+len(entry.prompt))

				// render the listbox widget
				listView := listBoxStyle.Render(entry.listModel.View())

				// write the first line of the listbox widget next to the prompt
				lines := strings.Split(listView, "\n")
				if len(lines) == 0 {
					log.Fatal("error while building the listbox. Exiting")
				}
				builder.WriteString(lines[0] + "\n")

				// add indentRemainingBoxLines to each remaining line of the list
				for _, line := range lines[1:] {
					builder.WriteString(indentRemainingBoxLines + line + "\n")
				}
			} else {
				// Show just label and selected value when not focused (aligned with text inputs)
				selectedDisplay := entry.selected
				if selectedDisplay == "" {
					selectedDisplay = "(none selected)"
				}
				builder.WriteString(fmt.Sprintf(" %s%s\n", entry.prompt, selectedDisplay))
			}
		} else {
			// Render text input
			if entry.textInput.Validate(entry.textInput.Value()) == nil {
				builder.WriteString(goodTick.Render("✔"))
			} else {
				builder.WriteString(badTick.Render("✖"))
			}
			builder.WriteString(fmt.Sprintf(" %s ", entry.textInput.View()))
			builder.WriteRune('\n')
		}
	}

	builder.WriteRune('\n')
	builder.WriteString(helpStyle.Render(" ↑, ↓: navigate, enter: confirm, ctrl+c: quit, ctrl+s: save"))
	if m.err != nil {
		builder.WriteString(errorStyle.Render(m.err.Error()))
	}
	return builder.String()
}

func (m *Model) Execute() error {
	if _, err := tea.NewProgram(m).Run(); err != nil {
		return err
	}
	return nil
}

func (m *Model) focusEntry(id int) {
	// Blur current entry
	if m.entries[m.focusIndex].inputType == InputTypeText {
		m.entries[m.focusIndex].textInput.Blur()
	}
	m.focusIndex = id
	// Focus new entry
	if m.entries[m.focusIndex].inputType == InputTypeText {
		m.entries[m.focusIndex].textInput.Focus()
	}
}

// nextEntryIfNoError focuses the next entry
func (m *Model) nextEntryIfNoError() {
	// Check for errors on current entry
	if m.entries[m.focusIndex].inputType == InputTypeText {
		m.err = m.entries[m.focusIndex].textInput.Err
		if m.err != nil {
			return
		}
	}

	// Blur current entry
	if m.entries[m.focusIndex].inputType == InputTypeText {
		m.entries[m.focusIndex].textInput.Blur()
	}

	m.focusIndex = (m.focusIndex + 1) % len(m.entries)

	// Focus new entry
	if m.entries[m.focusIndex].inputType == InputTypeText {
		m.entries[m.focusIndex].textInput.Focus()
	}
}

// prevEntryIfNoError focuses the previous entry
func (m *Model) prevEntryIfNoError() {
	// Check for errors on current entry
	if m.entries[m.focusIndex].inputType == InputTypeText {
		m.err = m.entries[m.focusIndex].textInput.Err
		if m.err != nil {
			return
		}
	}

	// Blur current entry
	if m.entries[m.focusIndex].inputType == InputTypeText {
		m.entries[m.focusIndex].textInput.Blur()
	}

	m.focusIndex--
	if m.focusIndex < 0 {
		m.focusIndex = len(m.entries) - 1
	}

	// Focus new entry
	if m.entries[m.focusIndex].inputType == InputTypeText {
		m.entries[m.focusIndex].textInput.Focus()
	}
}

func (m *Model) GetVariablesForInteractiveForm(iForm InteractiveForm) map[string]string {
	variables := make(map[string]string)
	for i := range iForm.Inputs {
		if m.entries[i].inputType == InputTypeList {
			variables[iForm.Inputs[i].Placeholder] = m.entries[i].selected
		} else {
			variables[iForm.Inputs[i].Placeholder] = m.entries[i].textInput.Value()
		}
	}
	return variables
}
