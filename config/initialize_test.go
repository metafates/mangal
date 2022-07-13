package config

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInitialize(t *testing.T) {
	Convey("Given that config is nil", t, func() {
		UserConfig = nil

		Convey("When initialize is called", func() {
			Initialize("", false)

			Convey("Then config should not be nil", func() {
				So(UserConfig, ShouldNotBeNil)

				Convey("And it should be valid", func() {
					err := ValidateConfig(UserConfig)
					So(err, ShouldBeNil)
				})
			})
		})
	})
}