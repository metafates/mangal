package chapters

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var _ help.KeyMap = (*KeyMap)(nil)

type KeyMap struct {
	UnselectAll,
	SelectAll,
	Toggle,
	Read,
	Download,
	Confirm key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Toggle,
		k.Read,
		k.Download,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
	}
}
