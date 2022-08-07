package mangadex

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var mangadex = New()

func TestMangadex_Search(t *testing.T) {
	Convey("Given a mangadex instance", t, func() {
		Convey("When searching for a manga", func() {
			mangas, err := mangadex.Search("Death Note")
			Convey("Then the error should be nil", func() {
				So(err, ShouldBeNil)

				Convey("And the result should be a list of mangas", func() {
					So(len(mangas), ShouldBeGreaterThan, 0)

					Convey("And each manga should have a name, URL and ID", func() {
						for _, manga := range mangas {
							So(manga.Name, ShouldNotBeEmpty)
							So(manga.URL, ShouldNotBeEmpty)
							So(manga.ID, ShouldNotBeEmpty)
						}
					})
				})
			})
		})
	})
}
