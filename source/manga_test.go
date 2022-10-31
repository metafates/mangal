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
	Name:     "Death Note",
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
				So(testManga.Dirname(), ShouldEqual, util.SanitizeFilename(testManga.Name))
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
						So(lo.Must(filesystem.Api().IsDir(path)), ShouldBeTrue)
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
						So(lo.Must(filesystem.Api().IsDir(path)), ShouldBeTrue)

						Convey("It should be a temp directory", func() {
							isTemp := lo.Must(filesystem.Api().Exists(filepath.Join(os.TempDir(), filepath.Base(path))))
							So(isTemp, ShouldBeTrue)
						})
					})
				})
			})
		})
	})
}

func TestManga_PopulateMetadata(t *testing.T) {
	Convey("Given a manga", t, func() {
		Convey("When PopulateMetadata is called", func() {
			err := testManga.PopulateMetadata(func(string) {})
			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)

				Convey("It should fetch the metadata", func() {
					So(testManga.Metadata.Cover, ShouldNotBeEmpty)
					So(testManga.Metadata.Summary, ShouldNotBeEmpty)
					So(len(testManga.Metadata.Genres), ShouldNotBeEmpty)
					So(len(testManga.Metadata.Tags), ShouldNotBeEmpty)
				})
			})
		})
	})
}

func TestManga_SeriesJSON(t *testing.T) {
	Convey("Given a manga", t, func() {
		Convey("When SeriesJSON is called", func() {
			buf := testManga.SeriesJSON()
			Convey("It should return a json buffer", func() {
				So(buf, ShouldNotBeEmpty)
			})
		})
	})
}
