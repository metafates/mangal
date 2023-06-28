package chapters

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/base"
)

var _ base.State = (*State)(nil)

type State struct {
	client   *libmangal.Client
	chapters []libmangal.Chapter
	list     list.Model
	keyMap   KeyMap
}

func (s *State) Intermediate() bool {
	return false
}

func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *State) Title() base.Title {
	item, ok := s.list.SelectedItem().(Item)
	if !ok {
		return base.Title{Text: "Chapters"}
	}

	volume := item.Volume()
	manga := volume.Manga()

	return base.Title{Text: fmt.Sprintf("%s / Vol. %d", manga.Info().Title, volume.Info().Number)}
}

func (s *State) Status() string {
	return ""
}

func (s *State) Backable() bool {
	return s.list.FilterState() == list.Unfiltered
}

func (s *State) Resize(size base.Size) {
	s.list.SetSize(size.Width, size.Height)
}

func (s *State) Update(model base.Model, msg tea.Msg) (cmd tea.Cmd) {
	s.list, cmd = s.list.Update(msg)
	return cmd
}

func (s *State) View(model base.Model) string {
	return s.list.View()
}

func (s *State) Init(model base.Model) tea.Cmd {
	return nil
}
