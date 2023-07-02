package textinput

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/mangalorg/mangal/tui/util"
)

func New(options Options) *State {
	if options.Title.Text == "" {
		options.Title.Text = "Search"
	}

	input := textinput.New()
	input.Placeholder = options.Placeholder

	return &State{
		options:   options,
		textinput: input,
		keyMap: KeyMap{
			Confirm: util.Bind("confirm", "enter"),
		},
	}
}
