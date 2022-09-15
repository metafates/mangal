package util

import (
	"github.com/metafates/mangal/filesystem"
	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
	"path/filepath"
	"testing"
)

func TestUnzip(t *testing.T) {
	Convey("Given a zip file", t, func() {
		filesystem.SetOsFs()
		path := filepath.Join(filepath.Dir(lo.Must(filepath.Abs("."))), filepath.Join("assets", "testdata", "zipdata.zip"))
		file := lo.Must(filesystem.Api().Open(path))
		filesystem.SetMemMapFs()
		Convey("When unzipping it", func() {
			err := Unzip(file, lo.Must(file.Stat()).Size(), "a")
			Convey("Then the error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("And the files should be extracted", func() {
					for _, info := range []lo.Tuple2[string, bool]{
						{filepath.Join("a", "zipdata", "hey.jpeg"), false},
						{filepath.Join("a", "zipdata", "a"), true},
						{filepath.Join("a", "zipdata", "a", "b"), true},
						{filepath.Join("a", "zipdata", "a", "hello.txt"), false},
					} {
						filename := info.A
						isDir := info.B

						exists := lo.Must(filesystem.Api().Exists(filename))
						So(exists, ShouldBeTrue)

						if isDir {
							isDir := lo.Must(filesystem.Api().IsDir(filename))
							So(isDir, ShouldBeTrue)
						}
					}
				})
			})
		})
	})
}
