package style

import "github.com/charmbracelet/lipgloss"

func New() lipgloss.Style {
	return lipgloss.NewStyle()
}

func NewColored(foreground, background lipgloss.Color) lipgloss.Style {
	return New().Foreground(foreground).Background(background)
}

func Fg(color lipgloss.Color) func(string) string {
	return NewColored(color, "").Render
}

func Bg(color lipgloss.Color) func(string) string {
	return NewColored("", color).Render
}

func Truncate(max int) func(string) string {
	return New().Width(max).Render
}

var (
	Faint     = New().Faint(true).Render
	Bold      = New().Bold(true).Render
	Italic    = New().Italic(true).Render
	Underline = New().Underline(true).Render
)
