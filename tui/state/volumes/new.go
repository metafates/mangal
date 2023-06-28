package volumes

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/util"
)

func New(client *libmangal.Client, volumes []libmangal.Volume) *State {
	return &State{
		client:  client,
		volumes: volumes,
		list: util.NewList(
			1,
			volumes,
			func(volume libmangal.Volume) list.Item {
				return Item{volume}
			},
		),
		keyMap: KeyMap{
			Confirm: util.Bind("confirm", "enter"),
		},
	}
}
