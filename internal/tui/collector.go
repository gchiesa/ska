package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type (
	errMsg error
)

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#416767"}
	highlight = lipgloss.AdaptiveColor{Light: "#83ADF4", Dark: "#83ADF4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	bad       = lipgloss.AdaptiveColor{Light: "#CE1E00", Dark: "#CE1E00"}

	headerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			PaddingLeft(2).PaddingRight(2).Foreground(special)

	focusedStyle = lipgloss.NewStyle().Bold(true).Foreground(highlight)
	blurredStyle = lipgloss.NewStyle().Bold(false).Foreground(subtle)
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle
	errorStyle   = lipgloss.NewStyle().Foreground(bad).MarginTop(2).MarginBottom(1)
)

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds = make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			var i int
			for i = range m.inputs {
				if err := m.inputs[i].Validate(m.inputs[i].Value()); err != nil {
					m.err = err
					break
				}
			}
			if m.err != nil {
				m.focusInput(i)
			} else {
				return m, tea.Quit
			}
		case tea.KeyEnter:
			if m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}
			m.nextInputIfNoError()
		case tea.KeyCtrlC, tea.KeyEsc:
			m.exitWithCtrlC = true
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP, tea.KeyUp:
			m.prevInputIfNoError()
		}

		for i := range m.inputs {
			m.inputs[i].PromptStyle = noStyle
			m.inputs[i].Blur()
		}
		m.inputs[m.focusIndex].PromptStyle = focusedStyle
		m.inputs[m.focusIndex].Focus()

	case errMsg:
		m.err = msg
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	var builder strings.Builder

	builder.WriteRune('\n')
	builder.WriteString(headerStyle.Render(m.header))
	builder.WriteRune('\n')
	for i := range m.inputs {
		builder.WriteRune('⇨')
		builder.WriteString(fmt.Sprintf(" %s ", m.inputs[i].View()))
		builder.WriteRune('\n')
	}
	builder.WriteRune('\n')
	builder.WriteString(helpStyle.Render(" ↑, ↓: navigate, enter: confirm, ctrl+c: quit, ctrl+s: save"))
	if m.err != nil {
		builder.WriteString(errorStyle.Render(m.err.Error()))
	}
	return builder.String()
}

func (m *Model) Execute() error {
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		return err
	}
	return nil
}

func (m *Model) focusInput(id int) {
	m.inputs[m.focusIndex].Blur()
	m.focusIndex = id
	m.inputs[m.focusIndex].Focus()
}

// nextInput focuses the next input field
func (m *Model) nextInputIfNoError() {
	m.err = m.inputs[m.focusIndex].Err
	if m.err != nil {
		return
	}
	m.inputs[m.focusIndex].Blur()
	m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
	m.inputs[m.focusIndex].Focus()
}

// prevInputIfNoError focuses the previous input field
func (m *Model) prevInputIfNoError() {
	m.err = m.inputs[m.focusIndex].Err
	if m.err != nil {
		return
	}
	m.inputs[m.focusIndex].Blur()
	m.focusIndex--
	// Wrap around
	if m.focusIndex < 0 {
		m.focusIndex = len(m.inputs) - 1
	}
	m.inputs[m.focusIndex].Focus()
}

func (m *Model) GetVariablesForInteractiveForm(iForm InteractiveForm) map[string]string {
	variables := make(map[string]string)
	for i := range iForm.Inputs {
		variables[iForm.Inputs[i].Placeholder] = m.inputs[i].Value()
	}
	return variables
}
