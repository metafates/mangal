package history

import (
	"fmt"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func init() {
	filesystem.SetMemMapFs()
}

func TestHistory(t *testing.T) {
	Convey("Given a chapter", t, func() {
		chapter := source.Chapter{
			Name:     "adwad",
			URL:      "dwaofa",
			Index:    42069,
			SourceID: "fwaiog",
			ID:       "fawfa",
			Pages:    nil,
		}
		manga := source.Manga{
			Name:     "dawf",
			URL:      "fwa",
			Index:    1337,
			SourceID: "sajfioaw",
			ID:       "wjakfkawgjj",
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
					So(chapters[fmt.Sprintf("%s (%s)", chapter.Manga.Name, chapter.SourceID)].Name, ShouldEqual, chapter.Name)
				})
			})
		})
	})
}
