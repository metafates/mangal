package providers

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/path"
	"github.com/mangalorg/mangal/tui/base"
	"github.com/mangalorg/mangal/tui/state/loading"
	"github.com/mangalorg/mangal/tui/state/mangas"
	"github.com/mangalorg/mangal/tui/state/textinput"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/bbolt"
	"github.com/philippgille/gokv/encoding"
	"github.com/pkg/errors"
	"net/http"
	"path/filepath"
	"time"
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
func (s *State) Status() string {
	return s.list.Paginator.View()
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
					newPersistentStore := func(name string) (gokv.Store, error) {
						dir := filepath.Join(path.CacheDir(), "anilist")
						if err := fs.FS.MkdirAll(dir, 0755); err != nil {
							return nil, err
						}

						return bbolt.NewStore(bbolt.Options{
							BucketName: name,
							Path:       filepath.Join(dir, name+".db"),
							Codec:      encoding.Gob,
						})
					}

					httpClient := &http.Client{
						Timeout: time.Minute,
					}

					anilistOptions := libmangal.DefaultAnilistOptions()

					var err error
					anilistOptions.QueryToIDsStore, err = newPersistentStore("query-to-id")
					if err != nil {
						return err
					}

					anilistOptions.IDToMangaStore, err = newPersistentStore("id-to-manga")
					if err != nil {
						return err
					}

					anilistOptions.TitleToIDStore, err = newPersistentStore("title-to-id")
					if err != nil {
						return err
					}

					anilistOptions.AccessTokenStore, err = newPersistentStore("access-token")
					if err != nil {
						return err
					}

					anilist := libmangal.NewAnilist(anilistOptions)

					options := libmangal.DefaultClientOptions()
					options.FS = fs.FS
					options.Anilist = &anilist
					options.HTTPClient = httpClient

					client, err := libmangal.NewClient(model.Context(), item, options)
					if err != nil {
						return err
					}

					return textinput.New(textinput.Options{
						Title:       "Search",
						Prompt:      "Search for a manga",
						Placeholder: "",
						OnResponse: func(response string) tea.Cmd {
							return tea.Sequence(
								func() tea.Msg {
									return loading.New("Searching")
								},
								func() tea.Msg {
									m, err := client.SearchMangas(model.Context(), response)
									if err != nil {
										return err
									}

									return mangas.New(client, m)
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
	s.list, cmd = s.list.Update(msg)
	return cmd
}

// View implements base.State.
func (s *State) View(model base.Model) string {
	return s.list.View()
}
