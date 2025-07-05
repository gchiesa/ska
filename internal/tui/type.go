package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	header        string
	showBanner    bool
	focusIndex    int
	inputs        []textinput.Model
	err           error
	exitWithCtrlC bool
}
