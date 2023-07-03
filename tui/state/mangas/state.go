package mangas

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/config"
	"github.com/mangalorg/mangal/tui/base"
	"github.com/mangalorg/mangal/tui/state/chapters"
	"github.com/mangalorg/mangal/tui/state/listwrapper"
	"github.com/mangalorg/mangal/tui/state/loading"
	"github.com/mangalorg/mangal/tui/state/volumes"
)

var _ base.State = (*State)(nil)

type State struct {
	client *libmangal.Client
	mangas []libmangal.Manga
	list   *listwrapper.State
	keyMap KeyMap
}

func (s *State) Intermediate() bool {
	return false
}

func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *State) Title() base.Title {
	return base.Title{Text: "Mangas"}
}

func (s *State) Subtitle() string {
	return s.list.Subtitle()
}

func (s *State) Status() string {
	return s.list.Status()
}

func (s *State) Backable() bool {
	return s.list.Backable()
}

func (s *State) Resize(size base.Size) {
	s.list.Resize(size)
}

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
		case key.Matches(msg, s.keyMap.Confirm):
			return tea.Sequence(
				func() tea.Msg {
					return loading.New("Loading...")
				},
				func() tea.Msg {
					v, err := s.client.MangaVolumes(model.Context(), item.Manga)
					if err != nil {
						return err
					}

					if len(v) != 1 || !config.Config.TUI.ExpandSingleVolume.Get() {
						return volumes.New(s.client, v)
					}

					volume := v[0]
					c, err := s.client.VolumeChapters(model.Context(), volume)
					if err != nil {
						return err
					}

					return chapters.New(s.client, volume, c)
				},
			)
		}
	}
end:
	return s.list.Update(model, msg)
}

func (s *State) View(model base.Model) string {
	return s.list.View(model)
}

func (s *State) Init(model base.Model) tea.Cmd {
	return s.list.Init(model)
}
