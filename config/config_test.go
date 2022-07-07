package config

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	Convey("Given a default config file", t, func() {
		Convey("When initialize is called", func() {
			Initialize("", false)

			Convey("Should have at least one source", func() {
				So(len(UserConfig.Scrapers), ShouldBeGreaterThan, 0)
			})
		})
	})
}

func TestParseConfig(t *testing.T) {
	Convey("Given a config file", t, func() {
		Convey("When parseConfig is called", func() {
			config, err := ParseConfig(DefaultConfigBytes)

			Convey("Then the error should be nil", func() {
				So(err, ShouldBeNil)

				Convey("And the config should be parsed", func() {
					So(config, ShouldNotBeNil)

					Convey("And the config should have at least one source", func() {
						So(len(config.Scrapers), ShouldBeGreaterThan, 0)
					})
				})
			})

		})
	})
}

func TestGetConfig(t *testing.T) {
	Convey("Given a config file", t, func() {
		Convey("When getConfig is called with empty string", func() {
			config := GetConfig("")

			Convey("Config should not be nil", func() {
				So(config, ShouldNotBeNil)
			})
		})
	})
}

func TestValidateConfig(t *testing.T) {
	Convey("Given a valid config file", t, func() {
		config := GetConfig("")
		Convey("When validateConfig is called", func() {
			err := ValidateConfig(config)

			Convey("Then the error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given an invalid config file", t, func() {
		config := GetConfig("")
		config.Scrapers = nil

		Convey("When validateConfig is called", func() {
			err := ValidateConfig(config)

			Convey("Then the error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
