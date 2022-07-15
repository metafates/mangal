package style

import "github.com/charmbracelet/lipgloss"

var (
	Common    = lipgloss.NewStyle().Margin(2, 2)
	Accent    = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	Secondary = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})
	Bold      = lipgloss.NewStyle().Bold(true)
	Italic    = lipgloss.NewStyle().Italic(true)
	Success   = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	Fail      = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	InputPrompt    = Accent.Copy().Bold(true)
	InputTitle     = InputPrompt.Copy()
	MangaListTitle = lipgloss.NewStyle().
			Background(lipgloss.Color("#9f86c0")).
			Foreground(lipgloss.Color("#231942")).
			Padding(0, 1)
	ChaptersListTitle = lipgloss.NewStyle().
				Background(lipgloss.Color("#e0b1cb")).
				Foreground(lipgloss.Color("#231942")).
				Padding(0, 1)
)
