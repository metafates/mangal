package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Back      key.Binding
	Quit      key.Binding
	Confirm   key.Binding
	Select    key.Binding
	SelectAll key.Binding

	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding

	ToStart key.Binding
	ToEnd   key.Binding

	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding

	state sessionState
}

func defaultKeyMap() keyMap {
	return keyMap{
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back")),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit")),
		Confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm")),
		Select: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "select")),
		ShowFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help")),
		CloseFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "close")),
	}
}

var stateKeyMap = map[sessionState]keyMap{
	searchState:         keyMapFor(searchState),
	spinnerState:        keyMapFor(spinnerState),
	mangaSelectState:    keyMapFor(mangaSelectState),
	chaptersSelectState: keyMapFor(chaptersSelectState),
	promptState:         keyMapFor(promptState),
	progressState:       keyMapFor(progressState),
	exitPromptState:     keyMapFor(exitPromptState),
}

// keyMapFor returns key map for specific state
func keyMapFor(state sessionState) keyMap {
	def := defaultKeyMap()
	nobind := key.NewBinding()
	def.state = state

	switch state {
	case searchState:
		def.Select = nobind
	case spinnerState:
		def.Select = nobind
		def.Confirm = nobind
	case mangaSelectState:
		def.Quit = key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("q", "quit"))
		addMoveBindings(&def)
	case chaptersSelectState:
		def.Quit = key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("q", "quit"))
		def.SelectAll = key.NewBinding(
			key.WithKeys("ctrl+a", "*"),
			key.WithHelp("ctrl+a/*", "select all"))
		addMoveBindings(&def)
	case promptState:
		def.Select = nobind
	case progressState:
		def.Back = nobind
		def.Confirm = nobind
		def.Select = nobind
	case exitPromptState:
		def.Select = nobind
		def.Quit = key.NewBinding(
			key.WithKeys("ctrl+c", "q", "enter"),
			key.WithHelp("q/enter", "quit"))
	}

	return def
}

func addMoveBindings(k *keyMap) {
	k.Up = key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"))

	k.Down = key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"))

	k.Left = key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "left"))

	k.Right = key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "right"))

	k.ToStart = key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "to start"))

	k.ToEnd = key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "to end"))
}
