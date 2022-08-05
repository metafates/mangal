package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (b *statefulBubble) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, b.loadSources())
}
