package providers

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/anilist"
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/tui/base"
	"github.com/mangalorg/mangal/tui/state/listwrapper"
	"github.com/mangalorg/mangal/tui/state/loading"
	"github.com/mangalorg/mangal/tui/state/mangas"
	"github.com/mangalorg/mangal/tui/state/textinput"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

var _ base.State = (*State)(nil)

type State struct {
	providersLoaders []libmangal.ProviderLoader
	list             *listwrapper.State
	keyMap           KeyMap
}

// Backable implements base.State.
func (s *State) Backable() bool {
	return s.list.Backable()
}

// Init implements base.State.
func (s *State) Init(model base.Model) tea.Cmd {
	return s.list.Init(model)
}

// Intermediate implements base.State.
func (s *State) Intermediate() bool {
	return s.list.Intermediate()
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

// Title implements base.State.
func (s *State) Title() base.Title {
	return base.Title{Text: "Providers"}
}

func (s *State) Subtitle() string {
	return s.list.Subtitle()
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
					return loading.New("Loading...", "")
				},
				func() tea.Msg {
					httpClient := &http.Client{
						Timeout: time.Minute,
					}

					options := libmangal.DefaultClientOptions()
					options.FS = fs.Afero
					options.Anilist = anilist.Client
					options.HTTPClient = httpClient

					client, err := libmangal.NewClient(model.Context(), item, options)
					if err != nil {
						return err
					}

					return textinput.New(textinput.Options{
						Title:  base.Title{Text: "Search"},
						Prompt: fmt.Sprintf("Using %q provider", client),
						OnResponse: func(response string) tea.Cmd {
							return tea.Sequence(
								func() tea.Msg {
									return loading.New("Loading", fmt.Sprintf("Searching for %q", response))
								},
								func() tea.Msg {
									m, err := client.SearchMangas(model.Context(), response)
									if err != nil {
										return err
									}

									return mangas.New(client, response, m)
								},
							)
						},
					})
				},
			)
		case key.Matches(msg, s.keyMap.info):
			return func() tea.Msg {
				return errors.New("not implemented")
			}
		}
	}
end:
	return s.list.Update(model, msg)
}

// View implements base.State.
func (s *State) View(model base.Model) string {
	return s.list.View(model)
}
