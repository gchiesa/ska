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
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
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

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}
			m.nextInputIfNoError()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP, tea.KeyUp:
			m.prevInputIfNoError()
		default:
			m.inputs[m.focusIndex].Err = nil
			m.inputs[m.focusIndex].Update(msg)
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
	//
	//for i := range m.inputs {
	//	m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	//}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var builder strings.Builder

	builder.WriteRune('\n')
	builder.WriteString(headerStyle.Render(m.header))
	builder.WriteRune('\n')
	for i := range m.inputs {
		builder.WriteRune('⇨')
		builder.WriteString(fmt.Sprintf(" %s ", m.inputs[i].View()))
		if i < len(m.inputs)-1 {
			builder.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	if m.err != nil {
		builder.WriteString(errorStyle.Render(fmt.Sprintf("%s", m.err.Error())))
	}

	fmt.Fprintf(&builder, "\n\n%s\n\n", *button)

	return builder.String()
}

func InputCollector(m tea.Model) error {
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		return err
	}
	return nil
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
