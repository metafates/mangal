package history

import (
	"fmt"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type testSource struct{}

func (testSource) Name() string {
	panic("")
}

func (testSource) Search(_ string) ([]*source.Manga, error) {
	panic("")
}

func (testSource) ChaptersOf(_ *source.Manga) ([]*source.Chapter, error) {
	panic("")
}

func (testSource) PagesOf(_ *source.Chapter) ([]*source.Page, error) {
	panic("")
}

func (testSource) ID() string {
	return "test source"
}

func init() {
	filesystem.SetMemMapFs()
}

func TestHistory(t *testing.T) {
	Convey("Given a chapter", t, func() {
		chapter := source.Chapter{
			Name:  "adwad",
			URL:   "dwaofa",
			Index: 42069,
			ID:    "fawfa",
			Pages: nil,
		}
		manga := source.Manga{
			Name:     "dawf",
			URL:      "fwa",
			Index:    1337,
			ID:       "wjakfkawgjj",
			Source:   testSource{},
			Chapters: []*source.Chapter{&chapter},
		}
		chapter.Manga = &manga

		Convey("When saving the chapter", func() {
			err := Save(&chapter)
			Convey("Then the error should be nil", func() {
				So(err, ShouldBeNil)

				Convey("And the chapter should be saved", func() {
					chapters, err := Get()
					So(err, ShouldBeNil)
					So(len(chapters), ShouldBeGreaterThan, 0)
					So(chapters[fmt.Sprintf("%s (%s)", chapter.Manga.Name, chapter.Source().ID())].Name, ShouldEqual, chapter.Name)
				})
			})
		})
	})
}
