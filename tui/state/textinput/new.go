package textinput

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/mangalorg/mangal/tui/util"
)

func New(options Options) *State {
	if options.Title == "" {
		options.Title = "Search"
	}

	return &State{
		options:   options,
		textinput: textinput.New(),
		keyMap: KeyMap{
			Confirm: util.Bind("confirm", "enter"),
		},
	}
}
