package providers

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/mangalorg/mangal/tui/state/listwrapper"
)

var _ help.KeyMap = (*KeyMap)(nil)

type KeyMap struct {
	info,
	confirm key.Binding
	list listwrapper.KeyMap
}

// FullHelp implements help.KeyMap.
func (p KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		p.ShortHelp(),
	}
}

// ShortHelp implements help.KeyMap.
func (p KeyMap) ShortHelp() []key.Binding {
	return append(
		p.list.ShortHelp(),
		p.confirm,
	)
}
