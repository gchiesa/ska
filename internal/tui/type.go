package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	header        string
	focusIndex    int
	inputs        []textinput.Model
	err           error
	exitWithCtrlC bool
}
