package errorstate

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var _ help.KeyMap = (*KeyMap)(nil)

type KeyMap struct{}

func (k KeyMap) ShortHelp() []key.Binding {
	return nil
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return nil
}
