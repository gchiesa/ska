package tui

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	header     string
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	err        error
}
