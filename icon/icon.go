package icon

import (
	"github.com/metafates/mangal/key"
	"github.com/spf13/viper"
)

const (
	emoji   = "emoji"
	nerd    = "nerd"
	plain   = "plain"
	kaomoji = "kaomoji"
	squares = "squares"
)

func AvailableVariants() []string {
	return []string{emoji, nerd, plain, kaomoji, squares}
}

type iconDef struct {
	emoji   string
	nerd    string
	plain   string
	kaomoji string
	squares string
}

func (i *iconDef) Get() string {
	switch viper.GetString(key.IconsVariant) {
	case emoji:
		return i.emoji
	case nerd:
		return i.nerd
	case plain:
		return i.plain
	case kaomoji:
		return i.kaomoji
	case squares:
		return i.squares
	default:
		return ""
	}
}

func Get(icon Icon) string {
	return icons[icon].Get()
}
