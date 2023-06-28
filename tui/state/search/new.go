package search

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/util"
)

func New(client *libmangal.Client) *State {
	return &State{
		client: client,
		keyMap: KeyMap{
			Confirm: util.Bind("confirm", "enter"),
		},
		textinput: textinput.New(),
	}
}
