package util

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewRwMap(t *testing.T) {
	Convey("When creating a new string rwmap", t, func() {
		rwmap := NewRwMap[string, string]()

		Convey("It should be empty", func() {
			So(rwmap.Len(), ShouldEqual, 0)
		})

		Convey("It should be the type of &RwMap[string, string]", func() {
			So(rwmap, ShouldHaveSameTypeAs, &RwMap[string, string]{})
		})
	})
}

func TestRwMap_Get(t *testing.T) {
	Convey("Given a non empty rwmap", t, func() {
		rwmap := NewRwMap[string, string]()
		rwmap.Set("key", "value")

		Convey("When getting a value that exists", func() {
			value, ok := rwmap.Get("key")

			Convey("It should indicate that it exists", func() {
				So(ok, ShouldBeTrue)
			})

			Convey("It should return the value", func() {
				So(value, ShouldEqual, "value")
			})
		})

		Convey("When getting a value that does not exist", func() {
			value, ok := rwmap.Get("non-existent")

			Convey("It should indicate that it does not exist", func() {
				So(ok, ShouldBeFalse)
			})

			Convey("It should return the default value", func() {
				So(value, ShouldEqual, "")
			})
		})
	})

	Convey("Given an empty rwmap", t, func() {
		rwmap := NewRwMap[string, string]()

		Convey("When getting a value that exists", func() {
			value, ok := rwmap.Get("key")

			Convey("It should indicate that it does not exist", func() {
				So(ok, ShouldBeFalse)
			})

			Convey("It should return the default value", func() {
				So(value, ShouldEqual, "")
			})
		})
	})
}

func TestRwMap_Set(t *testing.T) {
	Convey("Given a non empty rwmap", t, func() {
		rwmap := NewRwMap[string, string]()
		rwmap.Set("key", "value")

		Convey("When setting a value that exists", func() {
			rwmap.Set("key", "new-value")

			Convey("It should update the value", func() {
				value, _ := rwmap.Get("key")
				So(value, ShouldEqual, "new-value")
			})
		})

		Convey("When setting a value that does not exist", func() {
			rwmap.Set("non-existent", "new-value")

			Convey("It should add the value", func() {
				value, _ := rwmap.Get("non-existent")
				So(value, ShouldEqual, "new-value")
			})
		})
	})
}

func TestRwMap_Len(t *testing.T) {
	Convey("Given a non empty rwmap", t, func() {
		rwmap := NewRwMap[string, string]()
		rwmap.Set("key", "value")

		Convey("When getting the length", func() {
			length := rwmap.Len()

			Convey("It should return the correct length", func() {
				So(length, ShouldEqual, 1)
			})
		})
	})

	Convey("Given an empty rwmap", t, func() {
		rwmap := NewRwMap[string, string]()

		Convey("When getting the length", func() {
			length := rwmap.Len()

			Convey("It should return the correct length", func() {
				So(length, ShouldEqual, 0)
			})
		})
	})
}

func TestRwMap_Reset(t *testing.T) {
	Convey("Given a non empty rwmap", t, func() {
		rwmap := NewRwMap[string, string]()
		rwmap.Set("key", "value")

		Convey("When resetting the rwmap", func() {
			rwmap.Reset()

			Convey("It should be empty", func() {
				So(rwmap.Len(), ShouldEqual, 0)
			})
		})
	})
}
