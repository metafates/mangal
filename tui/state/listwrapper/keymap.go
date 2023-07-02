package listwrapper

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

var _ help.KeyMap = (*KeyMap)(nil)

type KeyMap struct {
	list *list.KeyMap
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.list.Filter,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
	}
}
