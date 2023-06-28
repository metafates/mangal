package providers

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/util"
)

func New(loaders []libmangal.ProviderLoader) *State {
	list_ := util.NewList(
		2,
		loaders,
		func(loader libmangal.ProviderLoader) list.Item {
			return Item{loader}
		},
	)

	return &State{
		providersLoaders: loaders,
		list:             list_,
		keyMap: KeyMap{
			info:    util.Bind("info", "i"),
			confirm: util.Bind("confirm", "enter"),
		},
	}
}
