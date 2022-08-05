package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Options struct {
	Continue bool
}

func Run(options *Options) error {

	bubble := newBubble()

	if options.Continue {
		bubble.newState(historyState)
	} else {
		bubble.newState(sourcesState)
	}

	if err := tea.NewProgram(bubble, tea.WithAltScreen()).Start(); err != nil {
		return err
	}

	return nil
}
