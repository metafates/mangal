package filesystem

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGet(t *testing.T) {
	Convey("When getting the filesystem", t, func() {
		fs := Get()
		Convey("Then the filesystem should be returned", func() {
			So(fs, ShouldNotBeNil)
		})
	})
}

func TestSetMemMapFs(t *testing.T) {
	Convey("When setting the filesystem to MemMapFS", t, func() {
		SetMemMapFs()
		Convey("Then the filesystem should be MemMapFS", func() {
			fs := Get()
			So(fs, ShouldNotBeNil)
			So(fs.Name(), ShouldEqual, "MemMapFS")
		})
	})
}

func TestSetOsFs(t *testing.T) {
	Convey("When setting the filesystem to OsFs", t, func() {
		SetOsFs()
		Convey("Then the filesystem should be OsFs", func() {
			fs := Get()
			So(fs, ShouldNotBeNil)
			So(fs.Name(), ShouldEqual, "OsFs")
		})
	})
}
