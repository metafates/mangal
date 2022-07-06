package style

import "github.com/charmbracelet/lipgloss"

var (
	CommonStyle         = lipgloss.NewStyle().Margin(2, 2)
	AccentStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BoldStyle           = lipgloss.NewStyle().Bold(true)
	InputPromptStyle    = AccentStyle.Copy().Bold(true)
	InputTitleStyle     = InputPromptStyle.Copy()
	SuccessStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	FailStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	MangaListTitleStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#9f86c0")).
				Foreground(lipgloss.Color("#231942")).
				Padding(0, 1)
	ChaptersListTitleStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#e0b1cb")).
				Foreground(lipgloss.Color("#231942")).
				Padding(0, 1)
)
