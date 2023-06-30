package confirm

import (
	"github.com/mangalorg/mangal/tui/util"
)

func New(message string, onResponse OnResponseFunc) *State {
	return &State{
		message: message,
		keyMap: KeyMap{
			Yes: util.Bind("yes", "y", "enter"),
			No:  util.Bind("no", "n"),
		},
		onResponse: onResponse,
	}
}
