package chapsdownloaded

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/mangal/tui/base"
	"github.com/mangalorg/mangal/tui/state/loading"
	"github.com/skratchdot/open-golang/open"
	"path/filepath"
)

var _ base.State = (*State)(nil)

type State struct {
	options Options
	keyMap  KeyMap
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
		case key.Matches(msg, s.keyMap.Open) && len(s.options.SucceedPaths) > 0:
			return tea.Sequence(
				func() tea.Msg {
					return loading.New("Opening")
				},
				func() tea.Msg {
					err := open.Run(filepath.Dir(s.options.SucceedPaths[0]))
					if err != nil {
						return err
					}

					return base.MsgBack{}
				},
			)
		case key.Matches(msg, s.keyMap.Retry) && len(s.options.Failed) > 0:
			return s.options.DownloadChapters(s.options.Failed)
		}
	}

	return nil
}

func (s *State) View(model base.Model) string {
	return fmt.Sprintf("%d succeed, %d failed", len(s.options.Succeed), len(s.options.Failed))
}

func (s *State) Init(model base.Model) tea.Cmd {
	return nil
}
