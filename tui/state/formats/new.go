package formats

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/state/listwrapper"
	"github.com/mangalorg/mangal/tui/util"
)

func New() *State {
	listWrapper := listwrapper.New(util.NewList(
		2,
		"manga", "mangas",
		libmangal.FormatValues(),
		func(format libmangal.Format) list.DefaultItem {
			return Item{format: format}
		},
	))

	return &State{
		list: listWrapper,
		keyMap: KeyMap{
			SetRead:     util.Bind("set for reading", "r"),
			SetDownload: util.Bind("set for downloading", "d"),
			SetAll:      util.Bind("set for all", "enter"),
			list:        listWrapper.GetKeyMap(),
		},
	}
}
