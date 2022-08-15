package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Options struct {
	Continue bool
	Install  bool
}

func Run(options *Options) error {

	bubble := newBubble()

	if options.Install {
		_, err := bubble.loadScrapers()
		if err != nil {
			return err
		}

		bubble.state = scrapersInstallState
	} else if options.Continue {
		_, err := bubble.loadHistory()
		if err != nil {
			return err
		}

		bubble.state = historyState
	} else {
		bubble.state = sourcesState
	}

	return tea.NewProgram(bubble, tea.WithAltScreen()).Start()
}
