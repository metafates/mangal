package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

func (k keyMap) ShortHelp() []key.Binding {
	return stateHelpMap[k.state].short
}

func (k keyMap) FullHelp() [][]key.Binding {
	return stateHelpMap[k.state].full
}

type helpBinds struct {
	short []key.Binding
	full  [][]key.Binding
}

var stateHelpMap = map[sessionState]helpBinds{
	searchState:         helpFor(searchState),
	mangaSelectState:    helpFor(mangaSelectState),
	chaptersSelectState: helpFor(chaptersSelectState),
	promptState:         helpFor(promptState),
	progressState:       helpFor(progressState),
	exitPromptState:     helpFor(exitPromptState),
}

// helpFor returns help binds for specific state.
func helpFor(state sessionState) helpBinds {
	var (
		short []key.Binding
		full  [][]key.Binding
	)

	k := stateKeyMap[state]

	switch state {
	case searchState:
		short = []key.Binding{k.Confirm, k.Quit}
		full = nil

	case spinnerState:
		short = []key.Binding{k.Back, k.Quit}
		full = nil

	case mangaSelectState:
		short = []key.Binding{k.Up, k.Down, k.Select, k.Back, k.Quit}
		full = [][]key.Binding{
			{k.Up, k.Down, k.Left, k.Right},
			{k.Select, k.Back, k.Quit}}

	case chaptersSelectState:
		short = []key.Binding{k.Up, k.Down, k.Select, k.SelectAll, k.Confirm, k.Back, k.Quit}
		full = [][]key.Binding{
			{k.Up, k.Down, k.Left, k.Right},
			{k.Select, k.SelectAll, k.Confirm, k.Back, k.Quit}}

	case promptState:
		short = []key.Binding{k.Confirm, k.Back, k.Quit}
		full = nil

	case progressState:
		short = []key.Binding{k.Quit}
		full = nil

	case exitPromptState:
		short = []key.Binding{k.Back, k.Quit}
		full = nil
	}

	return helpBinds{
		short: short,
		full:  full,
	}
}
