package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/mangal/tui/base"
	"github.com/mangalorg/mangal/tui/model"
)

func Run(state base.State) error {
	program := tea.NewProgram(model.New(state), tea.WithAltScreen())

	_, err := program.Run()
	return err
}
