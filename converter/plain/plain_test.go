package plain

import (
	"bytes"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func init() {
	lo.Must0(config.Setup())
	filesystem.SetMemMapFs()
}

func Test(t *testing.T) {
	plain := New()

	Convey("Given a plain converter", t, func() {
		Convey("When saving a chapter", func() {
			chapter := SampleChapter(t)
			result, err := plain.Save(chapter)
			Convey("Then the error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("And the result should be a path pointing to a directory", func() {
					So(result, ShouldNotBeEmpty)
					isDir, err := filesystem.Api().IsDir(result)

					if err != nil {
						t.Fatal(err)
					}

					So(isDir, ShouldBeTrue)

					Convey("And the directory should contain the chapter's pages", func() {
						files, err := filesystem.Api().ReadDir(result)
						if err != nil {
							t.Fatal(err)
						}

						So(len(files), ShouldEqual, len(chapter.Pages))

						lo.ForEach(files, func(file os.FileInfo, _ int) {
							So(file.Size(), ShouldBeGreaterThan, 0)
							So(file.IsDir(), ShouldBeFalse)
						})
					})
				})
			})
		})
	})
}
func SampleChapter(t *testing.T) *source.Chapter {
	t.Helper()
	chapter := source.Chapter{
		Name:  "chapter name",
		URL:   "chapter url",
		Index: 42069,
		ID:    "fawfa",
		Pages: []*source.Page{},
	}
	manga := source.Manga{
		Name:     "manga name",
		URL:      "manga url",
		Index:    1337,
		ID:       "wjakfkawgjj",
		Chapters: []*source.Chapter{&chapter},
	}
	chapter.Manga = &manga

	// to get images
	filesystem.SetOsFs()
	defer filesystem.SetMemMapFs()

	// get all images from ../assets/testdata
	err := filesystem.Api().Walk(
		// ../../assets/testdata
		// I wish windows used a normal path separator instead of whatever this \ is
		filepath.Join(filepath.Dir(filepath.Dir(lo.Must(filepath.Abs(".")))), filepath.Join("assets", "testdata")),
		func(path string, info fs.FileInfo, _ error) error {
			if lo.Must(filesystem.Api().IsDir(path)) || filepath.Ext(path) != ".jpg" {
				return nil
			}

			image, err := filesystem.Api().ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}

			page := source.Page{
				URL:       "dwadwaf",
				Index:     0,
				Extension: filepath.Ext(path),
				Chapter:   &chapter,
				Contents:  bytes.NewBuffer(image),
			}
			chapter.Pages = append(chapter.Pages, &page)

			return nil
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	return &chapter
}
