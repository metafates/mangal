package chapters

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/util"
)

func New(client *libmangal.Client, chapters []libmangal.Chapter) *State {
	return &State{
		client:   client,
		chapters: chapters,
		list: util.NewList(
			2,
			chapters,
			func(chapter libmangal.Chapter) list.Item {
				return Item{chapter}
			},
		),
		keyMap: KeyMap{},
	}
}
