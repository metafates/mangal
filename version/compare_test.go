package version

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCompareVersions(t *testing.T) {
	Convey("Given two versions with different patches", t, func() {
		v1, v2 := "1.0.0", "1.0.1"
		Convey("When comparing "+v1+" to "+v2, func() {
			result, err := Compare(v1, v2)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("Then the result should be -1", func() {
					So(result, ShouldEqual, -1)
				})
			})
		})

		Convey("When comparing "+v2+" to "+v1, func() {
			result, err := Compare(v2, v1)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("Then the result should be 1", func() {
					So(result, ShouldEqual, 1)
				})
			})
		})
	})

	Convey("Given two versions with different minor versions", t, func() {
		v1, v2 := "1.0.0", "1.1.0"
		Convey("When comparing "+v1+" to "+v2, func() {
			result, err := Compare(v1, v2)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("Then the result should be -1", func() {
					So(result, ShouldEqual, -1)
				})
			})
		})

		Convey("When comparing "+v2+" to "+v1, func() {
			result, err := Compare(v2, v1)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("Then the result should be 1", func() {
					So(result, ShouldEqual, 1)
				})
			})
		})
	})

	Convey("Given two versions with different major versions", t, func() {
		v1, v2 := "1.0.0", "2.0.0"
		Convey("When comparing "+v1+" to "+v2, func() {
			result, err := Compare(v1, v2)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("Then the result should be -1", func() {
					So(result, ShouldEqual, -1)
				})
			})
		})

		Convey("When comparing "+v2+" to "+v1, func() {
			result, err := Compare(v2, v1)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("Then the result should be 1", func() {
					So(result, ShouldEqual, 1)
				})
			})
		})
	})

	Convey("Given two same versions", t, func() {
		v1, v2 := "1.0.0", "1.0.0"
		Convey("When comparing "+v1+" to "+v2, func() {
			result, err := Compare(v1, v2)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("Then the result should be 0", func() {
					So(result, ShouldEqual, 0)
				})
			})
		})
	})
}
