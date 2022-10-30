package util

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPadZero(t *testing.T) {
	num := "123"
	Convey("Given a string "+num, t, func() {
		Convey("When padding with 4 zeros", func() {
			result := PadZero(num, 4)
			Convey("Then the result should be 0123", func() {
				So(result, ShouldEqual, "0123")
			})
		})

		Convey("When padding with 3 zeros", func() {
			result := PadZero(num, 3)
			Convey("Then the result should be 123", func() {
				So(result, ShouldEqual, "123")
			})
		})

		Convey("When padding with 2 zeros", func() {
			result := PadZero(num, 2)
			Convey("Then the result should be 123", func() {
				So(result, ShouldEqual, "123")
			})
		})

		Convey("When negative padding is performed", func() {
			result := PadZero(num, -1)
			Convey("Then the result should be 123", func() {
				So(result, ShouldEqual, "123")
			})
		})
	})
}

func TestFileStem(t *testing.T) {
	Convey("When the file name is 'foo.bar'", t, func() {
		result := FileStem("foo.bar")
		Convey("Then the result should be 'foo'", func() {
			So(result, ShouldEqual, "foo")
		})
	})
	Convey("When the file name is 'foo'", t, func() {
		result := FileStem("foo")
		Convey("Then the result should be 'foo'", func() {
			So(result, ShouldEqual, "foo")
		})
	})
	Convey("When the file name is 'foo.bar.baz'", t, func() {
		result := FileStem("foo.bar.baz")
		Convey("Then the result should be 'foo.bar'", func() {
			So(result, ShouldEqual, "foo.bar")
		})
	})
}

func TestQuantity(t *testing.T) {
	var (
		singular = "singular"
		plural   = "plural"
	)

	Convey("Given a quantity of 1", t, func() {
		quantity := 1
		Convey("When the quantity is converted to a string", func() {
			result := Quantify(quantity, singular, plural)
			Convey("Then the result should be '1 singular'", func() {
				So(result, ShouldEqual, "1 "+singular)
			})
		})
	})

	Convey("Given a quantity of 2", t, func() {
		quantity := 2
		Convey("When the quantity is converted to a string", func() {
			result := Quantify(quantity, singular, plural)
			Convey("Then the result should be '2 plural'", func() {
				So(result, ShouldEqual, "2 "+plural)
			})
		})
	})
}

func TestSanitizeFilename(t *testing.T) {
	invalidFilename := "~C:invalid/file name.txt."
	Convey("Given a string "+invalidFilename, t, func() {
		Convey("When the string is sanitized", func() {
			result := SanitizeFilename(invalidFilename)
			Convey("Then the result should be 'C_invalid_file_name.txt'", func() {
				So(result, ShouldEqual, "C_invalid_file_name.txt")
			})
		})
	})

	validFilename := "valid-file-name.txt"
	Convey("Given a string "+validFilename, t, func() {
		Convey("When the string is sanitized", func() {
			result := SanitizeFilename(validFilename)
			Convey("Then the result should be 'valid-file-name.txt'", func() {
				So(result, ShouldEqual, validFilename)
			})
		})
	})
}

func TestTerminalSize(t *testing.T) {
	t.Skipf("Cannot test terminal size")
}
