package formats

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/mangalorg/mangal/tui/state/listwrapper"
)

var _ help.KeyMap = (*KeyMap)(nil)

type KeyMap struct {
	SetRead, SetDownload key.Binding

	list listwrapper.KeyMap
}

// FullHelp implements help.KeyMap.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
	}
}

// ShortHelp implements help.KeyMap.
func (k KeyMap) ShortHelp() []key.Binding {
	return append(
		k.list.ShortHelp(),
		k.SetRead,
		k.SetDownload,
	)
}
