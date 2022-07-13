package tui

import (
	"github.com/metafates/mangal/config"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewBubble(t *testing.T) {
	config.Initialize("", false)

	Convey("When NewBubble is called", t, func() {
		bubble := NewBubble(SearchState)
		Convey("Then the bubble should not be nil", func() {
			So(bubble, ShouldNotBeNil)

			Convey("And the bubble should have an inital state", func() {
				So(bubble.state, ShouldEqual, SearchState)
			})

			Convey("And all components should exist", func() {
				So(bubble.input, ShouldNotBeNil)
				So(bubble.mangaList, ShouldNotBeNil)
				So(bubble.chaptersList, ShouldNotBeNil)
				So(bubble.ResumeList, ShouldNotBeNil)
				So(bubble.spinner, ShouldNotBeNil)
				So(bubble.help, ShouldNotBeNil)
				So(bubble.progress, ShouldNotBeNil)

				Convey("And the input should be empty", func() {
					So(bubble.input.Value(), ShouldEqual, "")

					Convey("Input also should be focused", func() {
						So(bubble.input.Focused(), ShouldBeTrue)
					})
				})
			})

			Convey("And channels should not be nil", func() {
				So(bubble.mangaChan, ShouldNotBeNil)
				So(bubble.chaptersChan, ShouldNotBeNil)
				So(bubble.chaptersProgressChan, ShouldNotBeNil)
				So(bubble.chapterPagesProgressChan, ShouldNotBeNil)
			})

		})
	})
}
