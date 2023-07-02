package chapsdownloaded

import (
	"github.com/mangalorg/mangal/tui/util"
)

func New(options Options) *State {
	state := &State{
		options: options,
	}

	state.keyMap = KeyMap{
		Open:  util.Bind("open directory", "o"),
		Quit:  util.Bind("quit", "q"),
		Retry: util.Bind("retry", "r"),
		state: state,
	}

	return state
}
