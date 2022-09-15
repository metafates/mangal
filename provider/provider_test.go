package provider

import (
	"github.com/metafates/mangal/provider/manganelo"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGet(t *testing.T) {
	Convey("When trying to get a valid provider", t, func() {
		_, ok := Get(manganelo.Config.Name)
		Convey("Then ok should be true", func() {
			So(ok, ShouldBeTrue)
		})
	})

	Convey("When trying to get an invalid provider", t, func() {
		_, ok := Get("kek")
		Convey("Then ok should be false", func() {
			So(ok, ShouldBeFalse)
		})
	})
}
