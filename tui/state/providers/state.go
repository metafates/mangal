package providers

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/base"
	"github.com/mangalorg/mangal/tui/state/loading"
	"github.com/mangalorg/mangal/tui/state/search"
	"github.com/pkg/errors"
)

var _ base.State = (*State)(nil)

type State struct {
	providersLoaders []libmangal.ProviderLoader
	list             list.Model
	keyMap           KeyMap
}

// Backable implements base.State.
func (s *State) Backable() bool {
	return s.list.FilterState() == list.Unfiltered
}

// Init implements base.State.
func (*State) Init(model base.Model) tea.Cmd {
	return nil
}

// Intermediate implements base.State.
func (*State) Intermediate() bool {
	return false
}

// KeyMap implements base.State.
func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

// Resize implements base.State.
func (s *State) Resize(size base.Size) {
	s.list.SetSize(size.Width, size.Height)
}

// Status implements base.State.
func (*State) Status() string {
	return ""
}

// Title implements base.State.
func (*State) Title() base.Title {
	return base.Title{Text: "Providers"}
}

// Update implements base.State.
func (s *State) Update(model base.Model, msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.FilterState() == list.Filtering {
			goto end
		}

		item, ok := s.list.SelectedItem().(Item)
		if !ok {
			return nil
		}

		switch {
		case key.Matches(msg, s.keyMap.confirm):
			return tea.Sequence(
				func() tea.Msg {
					return loading.New("Loading...")
				},
				func() tea.Msg {
					client, err := libmangal.NewClient(model.Context(), item, libmangal.DefaultClientOptions())
					if err != nil {
						return err
					}

					return search.New(client)
				},
			)
		case key.Matches(msg, s.keyMap.info):
			return func() tea.Msg {
				return errors.New("not implemented")
			}
		}
	}
end:
	s.list, cmd = s.list.Update(msg)
	return cmd
}

// View implements base.State.
func (s *State) View(model base.Model) string {
	return s.list.View()
}
