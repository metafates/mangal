package downloader

import (
	"bytes"
	"github.com/metafates/mangal/common"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/filesystem"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/afero"
	"log"
	"testing"
)

func init() {
	filesystem.Set(afero.NewMemMapFs())
	config.Initialize("", false)
}

func TestSaveTemp(t *testing.T) {
	Convey("Given a bytes buffer", t, func() {
		buffer := bytes.NewBuffer([]byte("test"))

		Convey("When saveTemp is called", func() {
			path, err := SaveTemp(buffer)

			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("It should return a path", func() {
				So(path, ShouldNotBeEmpty)

				Convey("That points to a non-empty file", func() {
					isEmpty, err := afero.IsEmpty(filesystem.Get(), path)

					So(err, ShouldBeNil)
					So(isEmpty, ShouldBeFalse)

					Convey("That contains the correct data", func() {
						contents, err := afero.ReadFile(filesystem.Get(), path)

						So(err, ShouldBeNil)
						So(contents, ShouldResemble, []byte("test"))
					})
				})
			})
		})
	})
}

func TestDownloadChapter(t *testing.T) {
	Convey("Given a chapter", t, func() {
		// get a chapter from the scraper
		manga, err := config.UserConfig.Scrapers[0].SearchManga(common.TestQuery)

		if err != nil {
			log.Fatal(err)
		}

		chapters, err := manga[0].Scraper.GetChapters(manga[0])

		if err != nil {
			log.Fatal(err)
		}

		chapter := chapters[0]

		Convey("When downloadChapter is called with PDF format", func() {
			config.UserConfig.Formats.Default = common.PDF

			path, err := DownloadChapter(chapter, nil, false)

			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("It should return a path", func() {
				So(path, ShouldNotBeEmpty)

				Convey("That points to a pdf file", func() {
					Convey("And it's not empty", func() {
						isEmpty, err := afero.IsEmpty(filesystem.Get(), path)

						So(err, ShouldBeNil)
						So(isEmpty, ShouldBeFalse)
					})
				})
			})
		})
	})
}
