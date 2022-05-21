package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

// Update handles all UI interactions.
func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch b.state {
	case searchState:
		return b.handleSearchState(msg)
	case spinnerState:
		return b.handleSpinnerState(msg)
	case mangaSelectState:
		return b.handleMangaSelectState(msg)
	case chaptersSelectState:
		return b.handleChaptersSelectState(msg)
	case promptState:
		return b.handlePromptState(msg)
	case progressState:
		return b.handleProgressState(msg)
	case exitPrompt:
		return b.handleExitPromptState(msg)
	}

	log.Fatal("Unknown state")
	return nil, nil
}
