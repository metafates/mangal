package converter

import (
	"github.com/metafates/mangal/constant"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGet(t *testing.T) {
	Convey("When trying to get a valid converter", t, func() {
		_, err := Get(constant.FormatCBZ)
		Convey("Then no error should be returned", func() {
			So(err, ShouldBeNil)
		})
	})

	Convey("When trying to get an invalid converter", t, func() {
		_, err := Get("kek")
		Convey("Then an error should be returned", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestAvailable(t *testing.T) {
	Convey("When getting the available converters", t, func() {
		converters := Available()
		Convey("Then the available converters should be returned", func() {
			So(converters, ShouldNotBeNil)
			So(len(converters), ShouldEqual, 4)
		})
	})
}
