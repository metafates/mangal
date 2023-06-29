package chapters

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/path"
	"github.com/mangalorg/mangal/tui/base"
	"github.com/pkg/errors"
	"github.com/zyedidia/generic/set"
)

var _ base.State = (*State)(nil)

type State struct {
	client   *libmangal.Client
	selected set.Set[*Item]
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
	item, ok := s.list.SelectedItem().(*Item)
	if !ok {
		return base.Title{Text: "Chapters"}
	}

	volume := item.chapter.Volume()
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
		case key.Matches(msg, s.keyMap.Download):
			options := libmangal.DownloadOptions{
				Format:              libmangal.FormatPDF,
				Directory:           ".",
				CreateMangaDir:      true,
				Strict:              false,
				SkipIfExists:        true,
				DownloadMangaCover:  false,
				DownloadMangaBanner: false,
				WriteSeriesJson:     false,
				WriteComicInfoXml:   false,
				ComicInfoXMLOptions: libmangal.DefaultComicInfoOptions(),
				ImageTransformer: func(bytes []byte) ([]byte, error) {
					return bytes, nil
				},
			}

			if s.selected.Size() == 0 {
				return downloadChapterCmd(
					model.Context(),
					s.client,
					item.chapter,
					options,
					func(path string) tea.Msg {
						// TODO: Return to some sort of 'download finished' screen
						return errors.New("unimplemented")
					},
				)
			}

			// unimplemented
			return nil
		case key.Matches(msg, s.keyMap.Read) || (s.selected.Size() == 0 && key.Matches(msg, s.keyMap.Confirm)):
			options := libmangal.DownloadOptions{
				Format:          libmangal.FormatPDF,
				Directory:       path.TempDir(),
				SkipIfExists:    true,
				ReadAfter:       true,
				ReadIncognito:   true,
				CreateMangaDir:  true,
				CreateVolumeDir: true,
				ImageTransformer: func(bytes []byte) ([]byte, error) {
					return bytes, nil
				},
			}

			return downloadChapterCmd(
				model.Context(),
				s.client,
				item.chapter,
				options,
				func(string) tea.Msg {
					return base.MsgBack{}
				},
			)
		case key.Matches(msg, s.keyMap.Toggle):

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
