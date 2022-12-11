package icon

import (
	"github.com/metafates/mangal/key"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"testing"
)

func TestGet(t *testing.T) {
	Convey("Given a icon", t, func() {
		i := Lua
		Convey("When getting the icon with emoji setting", func() {
			viper.Set(key.IconsVariant, emoji)
			result := Get(i)
			Convey("Then the result should be emoji icon", func() {
				So(result, ShouldEqual, icons[i].emoji)
			})
		})

		Convey("When getting the icon with nerd setting", func() {
			viper.Set(key.IconsVariant, nerd)
			result := Get(i)
			Convey("Then the result should be nerd icon", func() {
				So(result, ShouldEqual, icons[i].nerd)
			})
		})

		Convey("When getting the icon with plain setting", func() {
			viper.Set(key.IconsVariant, plain)
			result := Get(i)
			Convey("Then the result should be plain icon", func() {
				So(result, ShouldEqual, icons[i].plain)
			})
		})

		Convey("When getting the icon with kaomoji setting", func() {
			viper.Set(key.IconsVariant, kaomoji)
			result := Get(i)
			Convey("Then the result should be kaomoji icon", func() {
				So(result, ShouldEqual, icons[i].kaomoji)
			})
		})

		Convey("When getting the icon with no setting", func() {
			viper.Set(key.IconsVariant, "")
			result := Get(i)
			Convey("Then the result should be empty icon", func() {
				So(result, ShouldBeEmpty)
			})
		})
	})
}
