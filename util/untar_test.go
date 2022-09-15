package util

import (
	"github.com/metafates/mangal/filesystem"
	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
	"path/filepath"
	"testing"
)

func TestUntarGZ(t *testing.T) {
	Convey("Given a tar.gz file", t, func() {
		filesystem.SetOsFs()
		path := filepath.Join(filepath.Dir(lo.Must(filepath.Abs("."))), filepath.Join("assets", "testdata", "tardata.tar.gz"))
		file := lo.Must(filesystem.Api().Open(path))
		filesystem.SetMemMapFs()
		Convey("When untarring it", func() {
			err := UntarGZ(file, ".")
			Convey("Then the error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("And the files should be extracted", func() {
					for _, info := range []lo.Tuple2[string, bool]{
						{filepath.Join("tardata", "hey.jpeg"), false},
						{filepath.Join("tardata", "a"), true},
						{filepath.Join("tardata", "a", "b"), true},
						{filepath.Join("tardata", "a", "hello.txt"), false},
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
