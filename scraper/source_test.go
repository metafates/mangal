package scraper

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestValidateSource(t *testing.T) {
	Convey("Given an empty source", t, func() {
		source := &Source{}

		Convey("When validateSource is called", func() {
			err := ValidateSource(source)

			Convey("Then the error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a valid source", t, func() {
		source := &Source{
			Name:             "test",
			Base:             "https://example.com",
			ChaptersBase:     "https://example.com/chapters",
			SearchTemplate:   "https://example.com/search?q=%s",
			MangaAnchor:      "a",
			MangaTitle:       "a",
			ChapterAnchor:    "a",
			ChapterTitle:     "a",
			ReaderPage:       "a",
			RandomDelayMs:    0,
			ChaptersReversed: false,
			WhitespaceEscape: "%20",
		}

		Convey("When validateSource is called", func() {
			err := ValidateSource(source)

			Convey("Then the error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
