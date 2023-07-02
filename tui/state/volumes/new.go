package volumes

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/state/listwrapper"
	"github.com/mangalorg/mangal/tui/util"
)

func New(client *libmangal.Client, volumes []libmangal.Volume) *State {
	listWrapper := listwrapper.New(util.NewList(
		1,
		"volume", "volumes",
		volumes,
		func(volume libmangal.Volume) list.DefaultItem {
			return Item{volume}
		},
	))

	return &State{
		client:  client,
		volumes: volumes,
		list:    listWrapper,
		keyMap: KeyMap{
			Confirm: util.Bind("confirm", "enter"),
		},
	}
}
