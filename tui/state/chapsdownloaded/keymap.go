package chapsdownloaded

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var _ help.KeyMap = (*KeyMap)(nil)

type KeyMap struct {
	Quit,
	Retry key.Binding

	state *State
}

func (k KeyMap) ShortHelp() []key.Binding {
	bindings := []key.Binding{
		k.Quit,
	}

	if len(k.state.failed) > 0 {
		bindings = append(bindings, k.Retry)
	}

	return bindings
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
	}
}
