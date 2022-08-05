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
		bubble.setState(historyState)
	} else {
		bubble.setState(searchState)
	}

	if err := tea.NewProgram(bubble, tea.WithAltScreen()).Start(); err != nil {
		return err
	}

	return nil
}
