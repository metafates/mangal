package history

import (
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/scraper"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/afero"
	"testing"
)

func init() {
	filesystem.Set(afero.NewMemMapFs())
	config.Initialize("", false)
}

func TestHistory(t *testing.T) {
	Convey("Given a sample chapter url", t, func() {
		chapter := &scraper.URL{
			Relation: &scraper.URL{
				Address: "https://example.com",
				Info:    "Test manga",
				Index:   0,
			},
			Scraper: &scraper.Scraper{Source: &scraper.Source{Name: "Test source"}},
			Address: "https://example.com",
			Info:    "Test chapter",
			Index:   0,
		}

		Convey("When writeHistory is called", func() {
			err := WriteHistory(chapter)

			Convey("Then the history should be written", func() {
				So(err, ShouldBeNil)

				Convey("And history file should contain the correct data", func() {
					history, err := ReadHistory()
					So(err, ShouldBeNil)
					So(len(history), ShouldEqual, 1)

					So(history["https://example.com"].Chapter.Address, ShouldEqual, "https://example.com")
					So(history["https://example.com"].Chapter.Info, ShouldEqual, "Test chapter")
					So(history["https://example.com"].Chapter.Index, ShouldEqual, 0)
					So(history["https://example.com"].Manga.Address, ShouldEqual, "https://example.com")
					So(history["https://example.com"].Manga.Info, ShouldEqual, "Test manga")
					So(history["https://example.com"].Manga.Index, ShouldEqual, 0)
					So(history["https://example.com"].Manga.Scraper.Source.Name, ShouldEqual, "Test source")
				})
			})
		})
	})
}
