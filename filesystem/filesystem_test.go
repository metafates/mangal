package filesystem

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestApi(t *testing.T) {
	Convey("When getting the filesystem API", t, func() {
		api := Api()
		Convey("Then the filesystem API should be returned", func() {
			So(api, ShouldNotBeNil)
		})
	})
}

func TestSetMemMapFs(t *testing.T) {
	Convey("When setting the filesystem to MemMapFS", t, func() {
		SetMemMapFs()
		Convey("Then the filesystem should be MemMapFS", func() {
			api := Api()
			So(api, ShouldNotBeNil)
			So(api.Name(), ShouldEqual, "MemMapFS")
		})
	})
}

func TestSetOsFs(t *testing.T) {
	Convey("When setting the filesystem to OsFs", t, func() {
		SetOsFs()
		Convey("Then the filesystem should be OsFs", func() {
			api := Api()
			So(api, ShouldNotBeNil)
			So(api.Name(), ShouldEqual, "OsFs")
		})
	})
}
