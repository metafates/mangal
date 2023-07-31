package chapters

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/config"
	"github.com/mangalorg/mangal/path"
	"github.com/mangalorg/mangal/stringutil"
	"github.com/mangalorg/mangal/tui/base"
	"github.com/mangalorg/mangal/tui/state/anilistmangas"
	"github.com/mangalorg/mangal/tui/state/confirm"
	"github.com/mangalorg/mangal/tui/state/listwrapper"
	"github.com/mangalorg/mangal/tui/state/loading"
	"github.com/zyedidia/generic/set"
	"golang.org/x/exp/slices"
	"time"
)

var _ base.State = (*State)(nil)

type State struct {
	client   *libmangal.Client
	volume   libmangal.Volume
	selected set.Set[*Item]
	list     *listwrapper.State
	keyMap   KeyMap
}

func (s *State) Intermediate() bool {
	return false
}

func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *State) Title() base.Title {
	volume := s.volume
	manga := volume.Manga()

	return base.Title{Text: fmt.Sprintf("%s / Vol. %d", manga, volume.Info().Number)}
}

func (s *State) Subtitle() string {
	if s.selected.Size() > 0 {
		return fmt.Sprint(s.list.Subtitle(), " ", s.selected.Size(), " selected")
	}

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

		item, ok := s.list.SelectedItem().(*Item)
		if !ok {
			return nil
		}

		switch {
		case key.Matches(msg, s.keyMap.Toggle):
			item.Toggle()

			return nil
		case key.Matches(msg, s.keyMap.UnselectAll):
			for _, item := range s.selected.Keys() {
				item.Toggle()
			}

			return nil
		case key.Matches(msg, s.keyMap.SelectAll):
			for _, listItem := range s.list.Items() {
				item, ok := listItem.(*Item)
				if !ok {
					continue
				}

				if !item.IsSelected() {
					item.Toggle()
				}
			}

			return nil
		case key.Matches(msg, s.keyMap.Download) || (s.selected.Size() > 0 && key.Matches(msg, s.keyMap.Confirm)):
			format, err := libmangal.FormatString(config.Config.Download.Format.Get())
			if err != nil {
				return func() tea.Msg {
					return err
				}
			}

			options := libmangal.DownloadOptions{
				Format:              format,
				Directory:           config.Config.Download.Path.Get(),
				CreateVolumeDir:     config.Config.Download.Volume.CreateDir.Get(),
				CreateMangaDir:      config.Config.Download.Manga.CreateDir.Get(),
				Strict:              config.Config.Download.Strict.Get(),
				SkipIfExists:        config.Config.Download.SkipIfExists.Get(),
				DownloadMangaCover:  config.Config.Download.Manga.Cover.Get(),
				DownloadMangaBanner: config.Config.Download.Manga.Banner.Get(),
				WriteSeriesJson:     config.Config.Download.Metadata.SeriesJSON.Get(),
				WriteComicInfoXml:   config.Config.Download.Metadata.ComicInfoXML.Get(),
				ComicInfoXMLOptions: libmangal.DefaultComicInfoOptions(),
				ImageTransformer: func(bytes []byte) ([]byte, error) {
					return bytes, nil
				},
			}

			var chapters []libmangal.Chapter

			if s.selected.Size() == 0 {
				chapters = append(chapters, item.chapter)
			} else {
				for _, item := range s.selected.Keys() {
					chapters = append(chapters, item.chapter)
				}
			}

			slices.SortFunc(chapters, func(a, b libmangal.Chapter) int {
				if a.Info().Number < b.Info().Number {
					return -1
				}

				return 1
			})

			return func() tea.Msg {
				return confirm.New(
					fmt.Sprint("Download ", stringutil.Quantify(len(chapters), "chapter", "chapters")),
					func(response bool) tea.Cmd {
						if !response {
							return base.Back
						}

						return s.downloadChaptersCmd(chapters, options)
					},
				)
			}
		case key.Matches(msg, s.keyMap.Read) || (s.selected.Size() == 0 && key.Matches(msg, s.keyMap.Confirm)):
			format, err := libmangal.FormatString(config.Config.Read.Format.Get())
			if err != nil {
				return func() tea.Msg {
					return err
				}
			}

			var directory string
			if config.Config.Read.DownloadOnRead.Get() {
				directory = config.Config.Download.Path.Get()
			} else {
				directory = path.TempDir()
			}

			options := libmangal.DownloadOptions{
				Format:          format,
				Directory:       directory,
				SkipIfExists:    true,
				ReadAfter:       true,
				ReadIncognito:   config.Config.Read.Incognito.Get(),
				CreateMangaDir:  true,
				CreateVolumeDir: true,
				ImageTransformer: func(bytes []byte) ([]byte, error) {
					return bytes, nil
				},
			}

			return s.downloadChapterCmd(
				model.Context(),
				item.chapter,
				options,
			)
		case key.Matches(msg, s.keyMap.Anilist):
			return tea.Sequence(
				func() tea.Msg {
					return loading.New("Loading...", "Getting Anilist Mangas")
				},
				func() tea.Msg {
					var mangas []libmangal.AnilistManga

					mangaTitle := item.chapter.Volume().Manga().Info().Title

					closest, ok, err := s.client.Anilist().FindClosestManga(model.Context(), mangaTitle)
					if err != nil {
						return err
					}

					if ok {
						mangas = append(mangas, closest)
					}

					mangaSearchResults, err := s.client.Anilist().SearchMangas(model.Context(), mangaTitle)
					if err != nil {
						return nil
					}

					for _, manga := range mangaSearchResults {
						if manga.ID == closest.ID {
							continue
						}

						mangas = append(mangas, manga)
					}

					return anilistmangas.New(
						s.client.Anilist(),
						mangas,
						func(response *libmangal.AnilistManga) tea.Cmd {
							return tea.Sequence(
								func() tea.Msg {
									err := s.client.Anilist().BindTitleWithID(mangaTitle, response.ID)
									if err != nil {
										return err
									}

									return base.MsgBack{}
								},
								s.list.Notify("Binded to "+response.Title.English, time.Second*3),
							)
						},
					)
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
