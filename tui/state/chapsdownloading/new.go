package chapsdownloading

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/mangalorg/libmangal"
)

func New(client *libmangal.Client, chapters []libmangal.Chapter, options libmangal.DownloadOptions) *State {
	return &State{
		client:   client,
		chapters: chapters,
		options:  options,
		message:  "Preparing...",
		progress: progress.New(),
		spinner:  spinner.New(),
		keyMap:   KeyMap{},
	}
}
