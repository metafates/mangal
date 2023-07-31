package mangas

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/state/listwrapper"
	"github.com/mangalorg/mangal/tui/util"
)

func New(client *libmangal.Client, query string, mangas []libmangal.Manga) *State {
	listWrapper := listwrapper.New(util.NewList(
		2,
		"manga", "mangas",
		mangas,
		func(manga libmangal.Manga) list.DefaultItem {
			return Item{manga}
		},
	))

	return &State{
		query:  query,
		client: client,
		mangas: mangas,
		list:   listWrapper,
		keyMap: KeyMap{
			Confirm: util.Bind("confirm", "enter"),
			list:    listWrapper.GetKeyMap(),
		},
	}
}
