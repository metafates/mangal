package color

import "github.com/charmbracelet/lipgloss"

var (
	Red    = New("1")
	Green  = New("2")
	Yellow = New("3")
	Blue   = New("4")
	Purple = New("5")
	Cyan   = New("6")
	White  = New("7")
	Black  = New("8")
)

var (
	HiRed    = New("9")
	HiGreen  = New("10")
	HiYellow = New("11")
	HiBlue   = New("12")
	HiPurple = New("13")
	HiCyan   = New("14")
	HiWhite  = New("15")
	HiBlack  = New("16")
)

func New(color string) lipgloss.Color {
	return lipgloss.Color(color)
}
