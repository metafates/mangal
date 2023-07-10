package chapters

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/mangalorg/mangal/tui/state/listwrapper"
)

var _ help.KeyMap = (*KeyMap)(nil)

type KeyMap struct {
	UnselectAll,
	SelectAll,
	Toggle,
	Read,
	Download,
	Anilist,
	Confirm key.Binding

	list listwrapper.KeyMap
}

func (k KeyMap) ShortHelp() []key.Binding {
	return append(
		k.list.ShortHelp(),
		k.Toggle,
		k.Read,
		k.Download,
	)
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
		{k.SelectAll, k.UnselectAll},
		{k.Anilist},
	}
}
