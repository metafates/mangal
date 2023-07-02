package confirm

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/mangal/icon"
	"github.com/mangalorg/mangal/tui/base"
)

var _ base.State = (*State)(nil)

type OnResponseFunc func(response bool) tea.Cmd

type State struct {
	message    string
	keyMap     KeyMap
	onResponse OnResponseFunc
}

func (s *State) Intermediate() bool {
	return true
}

func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *State) Title() base.Title {
	return base.Title{Text: "Confirm"}
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
		case key.Matches(msg, s.keyMap.Yes):
			return s.onResponse(true)
		case key.Matches(msg, s.keyMap.No):
			return s.onResponse(false)
		}
	}

	return nil
}

func (s *State) View(model base.Model) string {
	return fmt.Sprintf("%s %s", icon.Question, s.message)
}

func (s *State) Init(model base.Model) tea.Cmd {
	return nil
}
