package chapsdownloading

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/mangalorg/libmangal"
)

func New(chapters []libmangal.Chapter, options Options) *State {
	return &State{
		options:  options,
		chapters: chapters,
		message:  "Preparing...",
		progress: progress.New(),
		spinner:  spinner.New(),
		keyMap:   KeyMap{},
	}
}
