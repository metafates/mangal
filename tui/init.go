package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the bubble
func (b Bubble) Init() tea.Cmd {
	return textinput.Blink
}
