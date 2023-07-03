package textinput

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/mangal/tui/base"
	"strings"
)

var _ base.State = (*State)(nil)

type OnResponseFunc func(response string) tea.Cmd

type State struct {
	options Options

	textinput textinput.Model
	keyMap    KeyMap
}

func (s *State) Intermediate() bool {
	return s.options.Intermediate
}

func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *State) Title() base.Title {
	return s.options.Title
}

func (s *State) Subtitle() string {
	return s.options.Prompt
}

func (s *State) Status() string {
	return ""
}

func (s *State) Backable() bool {
	return true
}

func (s *State) Resize(size base.Size) {
	s.textinput.Width = size.Width
}

func (s *State) Update(model base.Model, msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keyMap.Confirm) && strings.TrimSpace(s.textinput.Value()) != "":
			return s.options.OnResponse(strings.TrimSpace(s.textinput.Value()))
		}
	}

	s.textinput, cmd = s.textinput.Update(msg)
	return cmd
}

func (s *State) View(model base.Model) string {
	return s.textinput.View()
}

func (s *State) Init(model base.Model) tea.Cmd {
	return tea.Batch(s.textinput.Focus(), textinput.Blink)
}
