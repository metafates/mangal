package style

import "github.com/charmbracelet/lipgloss"

var (
	Bold      = lipgloss.NewStyle().Bold(true).Render
	Italic    = lipgloss.NewStyle().Italic(true).Render
	Underline = lipgloss.NewStyle().Underline(true).Render
)

var (
	Red     = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render
	Green   = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render
	Yellow  = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render
	Blue    = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Render
	Magenta = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render
	Cyan    = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render
	White   = lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render
	Black   = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render
)

var (
	HiBlack   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render
	HiRed     = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render
	HiGreen   = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render
	HiYellow  = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Render
	HiBlue    = lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Render
	HiMagenta = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Render
	HiCyan    = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render
	HiWhite   = lipgloss.NewStyle().Foreground(lipgloss.Color("16")).Render
)

func Hex(hex string) func(string) string {
	return func(s string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(hex)).Render(s)
	}
}

func Combined(styles ...func(string) string) func(string) string {
	return func(s string) string {
		for _, style := range styles {
			s = style(s)
		}
		return s
	}
}
