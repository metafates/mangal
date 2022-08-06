package style

import "github.com/charmbracelet/lipgloss"

type Style func(string) string

func Combined(styles ...func(string) string) func(string) string {
	return func(s string) string {
		for _, style := range styles {
			s = style(s)
		}
		return s
	}
}

func Padding(padding ...int) Style {
	return func(s string) string {
		return lipgloss.NewStyle().Padding(padding...).Render(s)
	}
}

func Trim(max int) Style {
	return func(s string) string {
		if len(s) <= max {
			return s
		}

		// Minus one for the ellipsis
		return s[:max-1] + "…"
	}
}
