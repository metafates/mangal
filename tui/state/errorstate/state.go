package errorstate

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/mangal/tui/base"
)

var _ base.State = (*State)(nil)

type State struct {
	error  error
	keyMap KeyMap
}

func (s *State) Intermediate() bool {
	return true
}

func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *State) Title() base.Title {
	// TODO: red bg
	return base.Title{Text: "Error"}
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
	return nil
}

func (s *State) View(model base.Model) string {
	return s.error.Error()
}

func (s *State) Init(model base.Model) tea.Cmd {
	return nil
}
