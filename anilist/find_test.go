package anilist

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFindClosest(t *testing.T) {
	Convey(`Given a query "Death Note"`, t, func() {
		query := "Death Note"
		Convey(`When trying to find closest on Anilist`, func() {
			result, err := FindClosest(query)
			Convey(`Then I should get a result`, func() {
				So(err, ShouldBeNil)
				So(result, ShouldNotBeEmpty)
				Convey(`And the first result should be "Death Note"`, func() {
					So(result.Title.English, ShouldEqual, "Death Note")
				})
			})
		})
	})
}
