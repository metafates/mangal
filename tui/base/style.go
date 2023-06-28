package base

import (
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Title,
	TitleBar,
	HelpBar lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Title: lipgloss.
			NewStyle().
			Bold(true).
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1),
		TitleBar: lipgloss.
			NewStyle().
			Padding(0, 0, 1, 2),
		HelpBar: lipgloss.
			NewStyle().
			Padding(0, 1),
	}
}
