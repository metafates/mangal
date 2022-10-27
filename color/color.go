package color

import "github.com/charmbracelet/lipgloss"

const (
	Red    = lipgloss.Color("1")
	Green  = lipgloss.Color("2")
	Yellow = lipgloss.Color("3")
	Blue   = lipgloss.Color("4")
	Purple = lipgloss.Color("5")
	Cyan   = lipgloss.Color("6")
	White  = lipgloss.Color("7")
	Black  = lipgloss.Color("8")
)

const (
	HiRed    = lipgloss.Color("9")
	HiGreen  = lipgloss.Color("10")
	HiYellow = lipgloss.Color("11")
	HiBlue   = lipgloss.Color("12")
	HiPurple = lipgloss.Color("13")
	HiCyan   = lipgloss.Color("14")
	HiWhite  = lipgloss.Color("15")
	HiBlack  = lipgloss.Color("16")
)

func New(s string) lipgloss.Color {
	return lipgloss.Color(s)
}
