package source

import (
	"fmt"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/key"
	"github.com/metafates/mangal/util"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"testing"
)

func init() {
	filesystem.SetMemMapFs()
	viper.Set(key.FormatsUse, constant.FormatPDF)
}

var testChapter = Chapter{
	Name:   "test chapter",
	URL:    "https://example.com",
	Index:  1,
	ID:     "test",
	Pages:  []*Page{},
	Manga:  &testManga,
	Volume: "1",
}

func TestChapter_Filename(t *testing.T) {
	Convey("Given a chapter", t, func() {
		Convey("When Filename is called", func() {
			Convey("It should return a sanitized filename", func() {
				const template = "&{index}! {chapter}// {volume} 28922@ {manga}"
				viper.Set(key.DownloaderChapterNameTemplate, template)
				filename := testChapter.Filename()

				Convey("It should match the given template", func() {
					So(filename, ShouldEqual, util.SanitizeFilename(fmt.Sprintf("&%d! %s// %s 28922@ %s.pdf", testChapter.Index, testChapter.Name, testChapter.Volume, testChapter.Manga.Name)))
				})
			})
		})
	})
}

func TestChapter_ComicInfoXML(t *testing.T) {
	Convey("Given a chapter", t, func() {
		Convey("When ComicInfo is called", func() {
			Convey("It should return a ComicInfo XML", func() {
				xml := testChapter.ComicInfo()
				So(xml, ShouldNotBeEmpty)
			})
		})
	})
}
