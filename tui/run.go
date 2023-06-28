package tui

import (
    "github.com/mangalorg/mangal/tui/model"
    tea "github.com/charmbracelet/bubbletea"
)

func Run() error {
    // TODO: add state
    program := tea.NewProgram(model.New(nil), tea.WithAltScreen())

    _, err := program.Run()
    return err
}
