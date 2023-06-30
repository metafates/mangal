package anilistmangas

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/base"
	"github.com/mangalorg/mangal/tui/state/loading"
	"github.com/mangalorg/mangal/tui/state/textinput"
)

var _ base.State = (*State)(nil)

type OnResponseFunc func(response *libmangal.AnilistManga) tea.Cmd

type State struct {
	anilist *libmangal.Anilist
	list    list.Model

	onResponse OnResponseFunc

	keyMap KeyMap
}

func (s *State) Intermediate() bool {
	return true
}

func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *State) Title() base.Title {
	return base.Title{Text: "Anilist Mangas"}
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.FilterState() == list.Filtering {
			goto end
		}

		switch {
		case key.Matches(msg, s.keyMap.Confirm):
			item, ok := s.list.SelectedItem().(Item)
			if !ok {
				return nil
			}

			return s.onResponse(item.Manga)
		case key.Matches(msg, s.keyMap.Search):
			return func() tea.Msg {
				return textinput.New(textinput.Options{
					Title:        "Search",
					Prompt:       "Enter Anilist manga title",
					Placeholder:  "",
					Intermediate: true,
					OnResponse: func(response string) tea.Cmd {
						return tea.Sequence(
							func() tea.Msg {
								return loading.New("Searching...")
							},
							func() tea.Msg {
								mangas, err := s.anilist.SearchMangas(model.Context(), response)
								if err != nil {
									return loading.New(err.Error())
								}

								return New(s.anilist, mangas, s.onResponse)
							},
						)
					},
				})
			}
		}
	}

end:
	s.list, cmd = s.list.Update(msg)
	return cmd
}

func (s *State) View(model base.Model) string {
	return s.list.View()
}

func (s *State) Init(model base.Model) tea.Cmd {
	return nil
}
