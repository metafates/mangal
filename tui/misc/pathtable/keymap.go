package pathtable

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var _ help.KeyMap = (*keyMap)(nil)

type keyMap struct {
	Copy,
	Quit key.Binding
}

// FullHelp implements help.KeyMap.
func (k *keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

// ShortHelp implements help.KeyMap.
func (k *keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Copy,
		k.Quit,
	}
}
