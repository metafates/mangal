package icon

import (
	"github.com/metafates/mangal/config"
	"github.com/spf13/viper"
)

const (
	emoji   = "emoji"
	nerd    = "nerd"
	plain   = "plain"
	kaomoji = "kaomoji"
)

func AvailableVariants() []string {
	return []string{emoji, nerd, plain, kaomoji}
}

type iconDef struct {
	emoji   string
	nerd    string
	plain   string
	kaomoji string
}

func (i *iconDef) Get() string {
	switch viper.GetString(config.IconsVariant) {
	case emoji:
		return i.emoji
	case nerd:
		return i.nerd
	case plain:
		return i.plain
	case kaomoji:
		return i.kaomoji
	default:
		return ""
	}
}

func Get(icon Icon) string {
	return icons[icon].Get()
}
