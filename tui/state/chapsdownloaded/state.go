package chapsdownloaded

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/base"
)

var _ base.State = (*State)(nil)

type State struct {
	client          *libmangal.Client
	options         libmangal.DownloadOptions
	succeed, failed []*libmangal.Chapter

	// FIXME: come up with a better solution
	// This is ugly, but avoids cyclic imports
	createChapsDownloadingState func(*libmangal.Client, []libmangal.Chapter, libmangal.DownloadOptions) base.State

	keyMap KeyMap
}

func (s *State) Intermediate() bool {
	return true
}

func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *State) Title() base.Title {
	return base.Title{Text: "Done"}
}

func (s *State) Subtitle() string {
	return ""
}

func (s *State) Status() string {
	return ""
}

func (s *State) Backable() bool {
	return true
}

func (s *State) Resize(size base.Size) {
}

func (s *State) Update(model base.Model, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keyMap.Quit):
			return tea.Quit
		case key.Matches(msg, s.keyMap.Retry):
			if len(s.failed) == 0 {
				return nil
			}

			var chapters = make([]libmangal.Chapter, len(s.failed))
			for i, chapter := range s.failed {
				chapters[i] = *chapter
			}

			return func() tea.Msg {
				return s.createChapsDownloadingState(
					s.client,
					chapters,
					s.options,
				)
			}
		}
	}

	return nil
}

func (s *State) View(model base.Model) string {
	return fmt.Sprintf("%d succeed, %d failed", len(s.succeed), len(s.failed))
}

func (s *State) Init(model base.Model) tea.Cmd {
	return nil
}
