package style

import "github.com/charmbracelet/lipgloss"

var (
	Common         = lipgloss.NewStyle().Margin(2, 2)
	Accent         = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	Bold           = lipgloss.NewStyle().Bold(true)
	InputPrompt    = Accent.Copy().Bold(true)
	InputTitle     = InputPrompt.Copy()
	Success        = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	Faile          = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	MangaListTitle = lipgloss.NewStyle().
			Background(lipgloss.Color("#9f86c0")).
			Foreground(lipgloss.Color("#231942")).
			Padding(0, 1)
	ChaptersListTitle = lipgloss.NewStyle().
				Background(lipgloss.Color("#e0b1cb")).
				Foreground(lipgloss.Color("#231942")).
				Padding(0, 1)
)
