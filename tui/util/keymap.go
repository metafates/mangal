package util

import "github.com/charmbracelet/bubbles/key"

func Bind(help string, primaryKey string, extraKeys ...string) key.Binding {
	var keys = make([]string, 1+len(extraKeys))
	keys[0] = primaryKey
	for i, k := range extraKeys {
		keys[i+1] = k
	}

	var primaryKeyHelp string
	if primaryKey == " " {
		primaryKeyHelp = "space"
	} else {
		primaryKeyHelp = primaryKey
	}

	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(primaryKeyHelp, help),
	)
}
