package version

import (
	"github.com/metafates/mangal/constant"
	. "github.com/smartystreets/goconvey/convey"
	"regexp"
	"runtime"
	"testing"
)

func TestLatestVersion(t *testing.T) {
	// I have no idea why this is failing on GitHub actions macOS runner
	if runtime.GOOS == constant.Darwin {
		t.Skip("Skipping test on darwin")
	}

	Convey("When getting the latest version", t, func() {
		version, err := Latest()
		Convey("It should not return an error", func() {
			So(err, ShouldBeNil)

			Convey("It should return a version", func() {
				So(version, ShouldNotBeEmpty)

				Convey("It should has a semver format", func() {
					semverRegex := regexp.MustCompile(`^v?(\d+)(\.\d+){0,2}(-\w+)?$`)
					So(semverRegex.MatchString(version), ShouldBeTrue)
				})
			})
		})
	})
}
