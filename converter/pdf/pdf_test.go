package pdf

import (
	"bytes"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/key"
	"github.com/metafates/mangal/source"
	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"io/fs"
	"path/filepath"
	"testing"
)

func init() {
	filesystem.SetMemMapFs()
	lo.Must0(config.Setup())
	viper.Set(key.FormatsUse, constant.FormatPDF)
}

func TestPDF(t *testing.T) {
	pdf := New()

	Convey("Given a FormatPDF converter", t, func() {
		Convey("When saving a chapter", func() {
			chapter := SampleChapter(t)
			result, err := pdf.Save(chapter)
			Convey("Then the error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("And the result should be a path with .pdf extension", func() {
					So(result, ShouldNotBeEmpty)
					So(filepath.Ext(result), ShouldEqual, ".pdf")

					Convey("A path that can be read", func() {
						file, err := filesystem.Api().Open(result)
						So(err, ShouldBeNil)
						So(file, ShouldNotBeNil)

						Convey("And the file should not be empty", func() {
							info := lo.Must(file.Stat())
							So(info.Size(), ShouldBeGreaterThan, 0)
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
			if lo.Must(filesystem.Api().IsDir(path)) || filepath.Ext(path) != ".jpeg" {
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
