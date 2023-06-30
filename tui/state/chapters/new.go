package chapters

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/util"
	"github.com/zyedidia/generic/set"
)

func New(client *libmangal.Client, chapters []libmangal.Chapter) *State {
	selectedSet := set.NewMapset[*Item]()
	return &State{
		client:   client,
		selected: selectedSet,
		list: util.NewList(
			2,
			chapters,
			func(chapter libmangal.Chapter) list.Item {
				return &Item{chapter: chapter, selectedItems: &selectedSet}
			},
		),
		keyMap: KeyMap{
			UnselectAll: util.Bind("unselect all", "backspace"),
			SelectAll:   util.Bind("select all", "a"),
			Toggle:      util.Bind("toggle", " "),
			Read:        util.Bind("read", "r"),
			Anilist:     util.Bind("anilist", "A"),
			Download:    util.Bind("download", "d"),
			Confirm:     util.Bind("confirm", "enter"),
		},
	}
}
