package source

import (
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"path/filepath"
	"testing"
)

func init() {
	filesystem.SetMemMapFs()
}

type testSource struct{}

func (t testSource) Name() string {
	return "test"
}

func (t testSource) ID() string {
	return "test"
}

func (t testSource) Search(string) (mangas []*Manga, err error) {
	return
}

func (t testSource) ChaptersOf(*Manga) (chapters []*Chapter, err error) {
	return
}

func (t testSource) PagesOf(*Chapter) (pages []*Page, err error) {
	return
}

var testManga = Manga{
	Name:     "test manga",
	URL:      "https://example.com",
	Index:    1,
	ID:       "test",
	Chapters: []*Chapter{},
	Source:   &testSource{},
}

func TestManga_Filename(t *testing.T) {
	Convey("Given a manga", t, func() {
		Convey("When Filename is called", func() {
			Convey("It should return a sanitized filename", func() {
				So(testManga.Filename(), ShouldEqual, util.SanitizeFilename(testManga.Name))
			})
		})
	})
}

func TestManga_Path(t *testing.T) {
	Convey("Given a manga", t, func() {
		Convey("When non-temp Path is called", func() {
			path, err := testManga.Path(false)
			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)

				Convey("It should return a path", func() {
					So(path, ShouldNotBeEmpty)

					Convey("It should be a directory", func() {
						So(lo.Must(filesystem.Get().IsDir(path)), ShouldBeTrue)
					})
				})
			})
		})

		Convey("When temp Path is called", func() {
			path, err := testManga.Path(true)
			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)

				Convey("It should return a path", func() {
					So(path, ShouldNotBeEmpty)

					Convey("It should be a directory", func() {
						So(lo.Must(filesystem.Get().IsDir(path)), ShouldBeTrue)

						Convey("It should be a temp directory", func() {
							isTemp := lo.Must(filesystem.Get().Exists(filepath.Join(os.TempDir(), filepath.Base(path))))
							So(isTemp, ShouldBeTrue)
						})
					})
				})
			})
		})
	})
}
