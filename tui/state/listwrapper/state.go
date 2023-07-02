package listwrapper

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/mangal/stringutil"
	"github.com/mangalorg/mangal/tui/base"
	"time"
)

var _ base.State = (*State)(nil)

type State struct {
	notification string
	list         list.Model
	keyMap       KeyMap
}

func (s *State) Intermediate() bool {
	return false
}

func (s *State) Backable() bool {
	return s.list.FilterState() == list.Unfiltered
}

func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *State) Title() base.Title {
	return base.Title{Text: "List"}
}

func (s *State) Subtitle() string {
	singular, plural := s.list.StatusBarItemName()
	subtitle := stringutil.Quantify(len(s.list.VisibleItems()), singular, plural)
	if s.list.FilterState() == list.FilterApplied {
		return fmt.Sprintf("%s %q", subtitle, s.list.FilterValue())
	}

	return subtitle
}

func (s *State) Status() string {
	if s.list.FilterState() == list.Filtering {
		return s.list.FilterInput.View()
	}

	if s.notification != "" {
		return s.list.Paginator.View() + " " + s.notification
	}

	return s.list.Paginator.View()
}

func (s *State) Resize(size base.Size) {
	s.list.SetSize(size.Width, size.Height)
}

func (s *State) Update(model base.Model, msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case NotificationMsg:
		s.notification = string(msg)
		return nil
	}

	s.list, cmd = s.list.Update(msg)
	return cmd
}

func (s *State) View(model base.Model) string {
	return s.list.View()
}

func (s *State) Init(model base.Model) tea.Cmd {
	return nil
}

func (s *State) FilterState() list.FilterState {
	return s.list.FilterState()
}

func (s *State) SelectedItem() list.Item {
	return s.list.SelectedItem()
}

func (s *State) GetKeyMap() KeyMap {
	return s.keyMap
}

func (s *State) Items() []list.Item {
	return s.list.Items()
}

func (s *State) Notify(message string, duration time.Duration) tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			return NotificationMsg(message)
		},
		func() tea.Msg {
			time.Sleep(duration)
			return NotificationMsg("")
		},
	)
}
