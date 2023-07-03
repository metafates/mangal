package chapsdownloading

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/color"
)

func New(chapters []libmangal.Chapter, options Options) *State {
	return &State{
		options:  options,
		chapters: chapters,
		message:  "Preparing...",
		progress: progress.New(),
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(color.Accent).Bold(true)),
		),
		keyMap: KeyMap{},
	}
}
