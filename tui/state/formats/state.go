package formats

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/mangal/tui/base"
	"github.com/mangalorg/mangal/tui/state/listwrapper"
)

var _ base.State = (*State)(nil)

type State struct {
	list   *listwrapper.State
	keyMap KeyMap
}

// Backable implements base.State.
func (*State) Backable() bool {
	return true
}

// Init implements base.State.
func (*State) Init(model base.Model) tea.Cmd {
	return nil
}

// Intermediate implements base.State.
func (*State) Intermediate() bool {
	return true
}

// KeyMap implements base.State.
func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

// Resize implements base.State.
func (s *State) Resize(size base.Size) {
	s.list.Resize(size)
}

// Status implements base.State.
func (s *State) Status() string {
	return s.list.Status()
}

// Subtitle implements base.State.
func (s *State) Subtitle() string {
	return s.list.Subtitle()
}

// Title implements base.State.
func (*State) Title() base.Title {
	return base.Title{Text: "Formats"}
}

// Update implements base.State.
func (s *State) Update(model base.Model, msg tea.Msg) tea.Cmd {
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
		case key.Matches(msg, s.keyMap.SetDownload):
			return func() tea.Msg {
				return item.SelectForDownloading()
			}
		case key.Matches(msg, s.keyMap.SetRead):
			return func() tea.Msg {
				return item.SelectForReading()
			}
		case key.Matches(msg, s.keyMap.SetAll):
			return tea.Batch(
				func() tea.Msg {
					return item.SelectForReading()
				},
				func() tea.Msg {
					return item.SelectForDownloading()
				},
			)
		}
	}
end:
	return s.list.Update(model, msg)
}

// View implements base.State.
func (s *State) View(model base.Model) string {
	return s.list.View(model)
}
