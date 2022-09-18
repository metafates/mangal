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

func TestWrap(t *testing.T) {
	s := "1234567890"
	Convey("Given a string "+s, t, func() {
		Convey("When wrapping with a width of 3", func() {
			expected := "123\n456\n789\n0"
			result := Wrap(s, 3)
			Convey("Then the result should be "+expected, func() {
				So(result, ShouldEqual, expected)
			})
		})

		Convey("When wrapping with a width of 2", func() {
			result := Wrap(s, 2)
			expected := "12\n34\n56\n78\n90"
			Convey("Then the result should be "+expected, func() {
				So(result, ShouldEqual, expected)
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
	plural := "Apples"
	Convey("Given a string "+plural, t, func() {
		Convey("When the quantity is 1", func() {
			result := Quantity(1, plural)
			Convey("Then the result should be '1 Apple'", func() {
				So(result, ShouldEqual, "1 Apple")
			})
		})
		Convey("When the quantity is 2", func() {
			result := Quantity(2, plural)
			Convey("Then the result should be '2 Apples'", func() {
				So(result, ShouldEqual, "2 Apples")
			})
			Convey("When the quantity is 0", func() {
				result := Quantity(0, plural)
				Convey("Then the result should be '0 Apples'", func() {
					So(result, ShouldEqual, "0 Apples")
				})
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

func TestCompareVersions(t *testing.T) {
	Convey("Given two versions with different patches", t, func() {
		v1, v2 := "1.0.0", "1.0.1"
		Convey("When comparing "+v1+" to "+v2, func() {
			result, err := CompareVersions(v1, v2)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("Then the result should be -1", func() {
					So(result, ShouldEqual, -1)
				})
			})
		})

		Convey("When comparing "+v2+" to "+v1, func() {
			result, err := CompareVersions(v2, v1)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("Then the result should be 1", func() {
					So(result, ShouldEqual, 1)
				})
			})
		})
	})

	Convey("Given two versions with different minor versions", t, func() {
		v1, v2 := "1.0.0", "1.1.0"
		Convey("When comparing "+v1+" to "+v2, func() {
			result, err := CompareVersions(v1, v2)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("Then the result should be -1", func() {
					So(result, ShouldEqual, -1)
				})
			})
		})

		Convey("When comparing "+v2+" to "+v1, func() {
			result, err := CompareVersions(v2, v1)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("Then the result should be 1", func() {
					So(result, ShouldEqual, 1)
				})
			})
		})
	})

	Convey("Given two versions with different major versions", t, func() {
		v1, v2 := "1.0.0", "2.0.0"
		Convey("When comparing "+v1+" to "+v2, func() {
			result, err := CompareVersions(v1, v2)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("Then the result should be -1", func() {
					So(result, ShouldEqual, -1)
				})
			})
		})

		Convey("When comparing "+v2+" to "+v1, func() {
			result, err := CompareVersions(v2, v1)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("Then the result should be 1", func() {
					So(result, ShouldEqual, 1)
				})
			})
		})
	})

	Convey("Given two same versions", t, func() {
		v1, v2 := "1.0.0", "1.0.0"
		Convey("When comparing "+v1+" to "+v2, func() {
			result, err := CompareVersions(v1, v2)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("Then the result should be 0", func() {
					So(result, ShouldEqual, 0)
				})
			})
		})
	})
}
