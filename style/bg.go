package style

import "github.com/charmbracelet/lipgloss"

var (
	BgRed     = lipgloss.NewStyle().Background(lipgloss.Color("1")).Render
	BgGreen   = lipgloss.NewStyle().Background(lipgloss.Color("2")).Render
	BgYellow  = lipgloss.NewStyle().Background(lipgloss.Color("3")).Render
	BgBlue    = lipgloss.NewStyle().Background(lipgloss.Color("4")).Render
	BgMagenta = lipgloss.NewStyle().Background(lipgloss.Color("5")).Render
	BgCyan    = lipgloss.NewStyle().Background(lipgloss.Color("6")).Render
	BgWhite   = lipgloss.NewStyle().Background(lipgloss.Color("7")).Render
	BgBlack   = lipgloss.NewStyle().Background(lipgloss.Color("8")).Render
)

var (
	BgHiBlack   = lipgloss.NewStyle().Background(lipgloss.Color("9")).Render
	BgHiRed     = lipgloss.NewStyle().Background(lipgloss.Color("10")).Render
	BgHiGreen   = lipgloss.NewStyle().Background(lipgloss.Color("11")).Render
	BgHiYellow  = lipgloss.NewStyle().Background(lipgloss.Color("12")).Render
	BgHiBlue    = lipgloss.NewStyle().Background(lipgloss.Color("13")).Render
	BgHiMagenta = lipgloss.NewStyle().Background(lipgloss.Color("14")).Render
	BgHiCyan    = lipgloss.NewStyle().Background(lipgloss.Color("15")).Render
	BgHiWhite   = lipgloss.NewStyle().Background(lipgloss.Color("16")).Render
)

func BgColor(color string) func(string) string {
	return func(s string) string {
		return lipgloss.NewStyle().Background(lipgloss.Color(color)).Render(s)
	}
}
