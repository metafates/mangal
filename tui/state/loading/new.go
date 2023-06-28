package loading

import (
	"github.com/charmbracelet/bubbles/spinner"
)

func New(message string) *State {
	return &State{
		message: message,
		spinner: spinner.New(spinner.WithSpinner(spinner.Dot)),
		keyMap:  KeyMap{},
	}
}
