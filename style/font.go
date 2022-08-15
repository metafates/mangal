package style

import "github.com/charmbracelet/lipgloss"

var (
	Bold      = lipgloss.NewStyle().Bold(true).Render
	Italic    = lipgloss.NewStyle().Italic(true).Render
	Underline = lipgloss.NewStyle().Underline(true).Render
	Faint     = lipgloss.NewStyle().Faint(true).Render
)
