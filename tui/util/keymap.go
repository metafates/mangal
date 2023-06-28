package util

import "github.com/charmbracelet/bubbles/key"

func Bind(help string, keys ...string) key.Binding {
	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], help),
	)
}
