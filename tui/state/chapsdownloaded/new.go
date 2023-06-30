package chapsdownloaded

import (
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/base"
	"github.com/mangalorg/mangal/tui/util"
)

func New(
	client *libmangal.Client,
	options libmangal.DownloadOptions,
	succeed,
	failed []*libmangal.Chapter,
	createChapsDownloadingState func(*libmangal.Client, []libmangal.Chapter, libmangal.DownloadOptions) base.State,
) *State {
	state := &State{
		client:                      client,
		options:                     options,
		succeed:                     succeed,
		failed:                      failed,
		createChapsDownloadingState: createChapsDownloadingState,
	}

	state.keyMap = KeyMap{
		Quit:  util.Bind("quit", "q"),
		Retry: util.Bind("retry", "r"),
		state: state,
	}

	return state
}
