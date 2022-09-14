package updater

import (
	. "github.com/smartystreets/goconvey/convey"
	"regexp"
	"testing"
)

func TestLatestVersion(t *testing.T) {
	Convey("When getting the latest version", t, func() {
		version, err := LatestVersion()
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
