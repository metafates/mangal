package model

import (
	"github.com/mangalorg/mangal/tui/util"
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Back, Quit, Help key.Binding
}

func newKeyMap() *keyMap {
	return &keyMap{
		Back: util.Bind("back", "esc"),
		Quit: util.Bind("quit", "ctrl+c", "ctrl+d"),
		Help: util.Bind("help", "?"),
	}
}
