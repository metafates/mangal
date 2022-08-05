package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (_ *statefulBubble) Init() tea.Cmd {
	return textinput.Blink
}
