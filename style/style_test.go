package style

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTrim(t *testing.T) {
	Convey("Given a string", t, func() {
		s := "lorem ipsum dolor sit amet"
		Convey("When trimming with a max of 10", func() {
			result := Truncate(10)(s)
			Convey("Then the result should be 'lorem ipsu…'", func() {
				So(result, ShouldEqual, "lorem ips…")
			})
		})

		Convey("When trimming with a max of 30h", func() {
			result := Truncate(30)(s)
			Convey("Then the result should be lorem ipsum dolor sit amet", func() {
				So(result, ShouldEqual, "lorem ipsum dolor sit amet")
			})
		})
	})
}

func TestCombined(t *testing.T) {
	Convey("Given a string", t, func() {
		s := "lorem ipsum dolor sit amet"
		Convey("When using combined with red and italic", func() {
			res := Combined(Red, Italic)(s)
			Convey("Then the result should be the same as Italic(Red(string))", func() {
				So(res, ShouldEqual, Italic(Red(s)))
			})
		})

		Convey("When using combined without arguments", func() {
			res := Combined()(s)
			Convey("Then the result should be the same as original", func() {
				So(res, ShouldEqual, s)
			})
		})
	})
}
