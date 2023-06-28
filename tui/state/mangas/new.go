package mangas

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/util"
)

func New(client *libmangal.Client, mangas []libmangal.Manga) *State {
	return &State{
		client: client,
		mangas: mangas,
		list: util.NewList(
			2,
			mangas,
			func(manga libmangal.Manga) list.Item {
				return Item{manga}
			},
		),
		keyMap: KeyMap{
			Confirm: util.Bind("confirm", "enter"),
		},
	}
}
