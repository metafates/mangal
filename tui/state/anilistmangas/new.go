package anilistmangas

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/util"
)

func New(anilist *libmangal.Anilist, chapters []libmangal.AnilistManga, onResponse OnResponseFunc) *State {
	return &State{
		anilist: anilist,
		list: util.NewList(2, chapters, func(manga libmangal.AnilistManga) list.Item {
			return Item{Manga: &manga}
		}),
		onResponse: onResponse,
		keyMap: KeyMap{
			Confirm: util.Bind("confirm", "enter"),
			Search:  util.Bind("search", "s"),
		},
	}
}
