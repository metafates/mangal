package loading

import (
	"github.com/charmbracelet/bubbles/spinner"
)

func New(message, subtitle string) *State {
	return &State{
		message:  message,
		subtitle: subtitle,
		spinner:  spinner.New(spinner.WithSpinner(spinner.Dot)),
		keyMap:   KeyMap{},
	}
}
