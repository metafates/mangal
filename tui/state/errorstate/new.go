package errorstate

import "github.com/mangalorg/mangal/tui/util"

func New(err error) *State {
	return &State{
		error: err,
		keyMap: KeyMap{
			Quit:      util.Bind("quit", "q"),
			CopyError: util.Bind("copy error", "c"),
		},
	}
}
