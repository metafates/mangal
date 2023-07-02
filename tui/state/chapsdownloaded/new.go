package chapsdownloaded

import (
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/base"
	"github.com/mangalorg/mangal/tui/util"
)

func New(
	client *libmangal.Client,
	options libmangal.DownloadOptions,
	dir string,
	succeed,
	failed []*libmangal.Chapter,
	createChapsDownloadingState func(*libmangal.Client, []libmangal.Chapter, libmangal.DownloadOptions) base.State,
) *State {
	state := &State{
		dir:                         dir,
		client:                      client,
		options:                     options,
		succeed:                     succeed,
		failed:                      failed,
		createChapsDownloadingState: createChapsDownloadingState,
	}

	state.keyMap = KeyMap{
		Open:  util.Bind("open directory", "o"),
		Quit:  util.Bind("quit", "q"),
		Retry: util.Bind("retry", "r"),
		state: state,
	}

	return state
}
