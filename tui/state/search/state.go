package search

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/base"
	"github.com/mangalorg/mangal/tui/state/loading"
	"github.com/mangalorg/mangal/tui/state/mangas"
)

var _ base.State = (*State)(nil)

type State struct {
	client    *libmangal.Client
	keyMap    KeyMap
	textinput textinput.Model
}

func (s *State) Intermediate() bool {
	return false
}

func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *State) Title() base.Title {
	return base.Title{Text: "Search"}
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
		case key.Matches(msg, s.keyMap.Confirm):
			return tea.Sequence(
				func() tea.Msg {
					return loading.New("Searching")
				},
				func() tea.Msg {
					m, err := s.client.SearchMangas(model.Context(), s.textinput.Value())
					if err != nil {
						return nil
					}

					return mangas.New(s.client, m)
				},
			)
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
