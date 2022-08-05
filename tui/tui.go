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
		bubble.state = historyState
	} else {
		bubble.state = sourcesState
	}

	if err := tea.NewProgram(bubble, tea.WithAltScreen()).Start(); err != nil {
		return err
	}

	return nil
}
