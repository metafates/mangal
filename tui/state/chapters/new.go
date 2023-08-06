package chapters

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/state/listwrapper"
	"github.com/mangalorg/mangal/tui/util"
	"github.com/zyedidia/generic/set"
)

func New(client *libmangal.Client, volume libmangal.Volume, chapters []libmangal.Chapter) *State {
	selectedSet := set.NewMapset[*Item]()
	listWrapper := listwrapper.New(util.NewList(
		2,
		"chapter", "chapters",
		chapters,
		func(chapter libmangal.Chapter) list.DefaultItem {
			return &Item{
				chapter:       chapter,
				selectedItems: &selectedSet,
				client:        client,
			}
		},
	))

	return &State{
		client:   client,
		volume:   volume,
		selected: selectedSet,
		list:     listWrapper,
		keyMap: KeyMap{
			UnselectAll: util.Bind("unselect all", "backspace"),
			SelectAll:   util.Bind("select all", "a"),
			Toggle:      util.Bind("toggle", " "),
			Read:        util.Bind("read", "r"),
			Anilist:     util.Bind("anilist", "A"),
			Download:    util.Bind("download", "d"),
			Confirm:     util.Bind("confirm", "enter"),
			list:        listWrapper.GetKeyMap(),
		},
	}
}
