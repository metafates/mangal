package util

import (
	"fmt"
	"github.com/metafates/mangal/filesystem"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"os"
	"regexp"
	"testing"
)

func TestPadZeros(t *testing.T) {
	Convey("Given some 2 digit integer", t, func() {
		const number = 42

		Convey("When I pad it with 5 zeros", func() {
			const padding = 5

			Convey("Then it should have 3 leading zeros", func() {
				So(PadZeros(number, padding), ShouldEqual, "00042")
			})
		})

		Convey("When I pad it with 3 zeros", func() {
			const padding = 3

			Convey("Then it should have 1 leading zero", func() {
				So(PadZeros(number, padding), ShouldEqual, "042")
			})
		})

		Convey("When I pad it with 1 or 2 zeros", func() {
			Convey("Then it should remain the same", func() {
				So(PadZeros(number, 1), ShouldEqual, "42")
				So(PadZeros(number, 2), ShouldEqual, "42")
			})
		})
	})
}

func TestPlural(t *testing.T) {
	Convey("Given a singular word", t, func() {
		const word = "book"

		Convey("When I pluralize it", func() {
			Convey("Then it should remain the same", func() {
				So(Plural(word, 2), ShouldEqual, "books")
			})
		})

		Convey("When I do not pluralize it", func() {
			Convey("Then it should be the same word", func() {
				So(Plural(word, 1), ShouldEqual, word)
			})
		})
	})

	Convey("Given a plural word", t, func() {
		const word = "books"

		Convey("When I pluralize it", func() {
			Convey("Then it should remain the same", func() {
				So(Plural(word, 2), ShouldEqual, "books")
			})
		})

		Convey("When I do not pluralize it", func() {
			Convey("Then it should be the same word", func() {
				So(Plural(word, 1), ShouldEqual, word)
			})
		})
	})
}

func TestMap(t *testing.T) {
	Convey("Given a list of integers", t, func() {
		var list = []int{1, 2, 3, 4, 5}

		Convey("When I map it with function that returns string", func() {
			Convey("Then it should return a list of strings", func() {
				So(Map(list, func(i int) string {
					return "x"
				}), ShouldResemble, []string{"x", "x", "x", "x", "x"})
			})
		})

		Convey("When I map it with function that does nothing", func() {
			Convey("Then it should return the same list", func() {
				So(Map(list, func(i int) int {
					return i
				}), ShouldResemble, []int{1, 2, 3, 4, 5})
			})
		})

		Convey("When I map it with function that returns nil", func() {
			Convey("Then it should return an empty list", func() {
				So(Map(list, func(i int) *int {
					return nil
				}), ShouldResemble, []*int{nil, nil, nil, nil, nil})
			})
		})
	})

	Convey("Given an empty list", t, func() {
		var list []int

		Convey("When I map it with function that returns string", func() {
			Convey("Then it should return an empty list", func() {
				So(Map(list, func(i int) string {
					return "x"
				}), ShouldResemble, []string{})
			})
		})
	})
}

func TestMax(t *testing.T) {
	Convey("Given two values", t, func() {
		const a = 1
		const b = 2

		Convey("When I get the max", func() {
			Convey("Then it should return the bigger one", func() {
				So(Max(a, b), ShouldEqual, b)
			})
		})
	})
}

func TestIsUnique(t *testing.T) {
	Convey("Given a list of unique elements", t, func() {
		var list = []int{1, 2, 3, 4, 5}

		Convey("When I check if it's unique", func() {
			Convey("Then it should return true", func() {
				So(IsUnique(list), ShouldBeTrue)
			})
		})
	})

	Convey("Given a list of non-unique elements", t, func() {
		var list = []int{1, 2, 3, 4, 5, 5}

		Convey("When I check if it's unique", func() {
			Convey("Then it should return false", func() {
				So(IsUnique(list), ShouldBeFalse)
			})
		})
	})
}

func TestSanitizeFilename(t *testing.T) {
	Convey("Given a filename that contains whitespaces", t, func() {
		const filename = "file name.ext"

		Convey("When I sanitize it", func() {
			Convey("Then whitespaces should be replaced with underscores", func() {
				So(SanitizeFilename(filename), ShouldEqual, "file_name.ext")
			})
		})
	})

	Convey("Given a filename that contains trailing dot", t, func() {
		const filename = "file.name."

		Convey("When I sanitize it", func() {
			Convey("Then trailing dot should be removed", func() {
				So(SanitizeFilename(filename), ShouldEqual, "file.name")
			})
		})
	})

	Convey("Given a filename that contains leading dot", t, func() {
		const filename = ".file.name"

		Convey("When I sanitize it", func() {
			Convey("Then leading dot should be removed", func() {
				So(SanitizeFilename(filename), ShouldEqual, "file.name")
			})
		})
	})

	Convey("Given a filename that contains leading and trailing dot", t, func() {
		const filename = ".file.name."

		Convey("When I sanitize it", func() {
			Convey("Then leading and trailing dot should be removed", func() {
				So(SanitizeFilename(filename), ShouldEqual, "file.name")
			})
		})
	})

	Convey("Given a filename that contains multiple whitespaces", t, func() {
		const filename = "file   name.ext"

		Convey("When I sanitize it", func() {
			Convey("Then whitespaces should be replaced with a single underscores", func() {
				So(SanitizeFilename(filename), ShouldEqual, "file_name.ext")
			})
		})
	})
}

func TestPrettyTrim(t *testing.T) {
	Convey("Given a long string", t, func() {
		const longString = "This is a very long string that should be trimmed"

		Convey("When I trim it", func() {
			Convey("Then it should be trimmed with trailing ellipse", func() {
				So(PrettyTrim(longString, 10), ShouldEqual, "This is...")
			})
		})
	})

	Convey("Given a short string", t, func() {
		const shortString = "This is a short string"

		Convey("When I trim it", func() {
			Convey("Then it should be the same string", func() {
				So(PrettyTrim(shortString, 30), ShouldEqual, shortString)
			})
		})
	})
}

func TestFind(t *testing.T) {
	Convey("Given a list of integers", t, func() {
		var list = []int{1, 2, 3, 4, 5}

		Convey("When I find an element", func() {
			element, ok := Find(list, func(i int) bool {
				return i == 3
			})

			Convey("Then it should indicate that it was found", func() {
				So(ok, ShouldBeTrue)
			})

			Convey("Then it should return the element", func() {
				So(element, ShouldEqual, 3)
			})
		})
	})

	Convey("Given a list of integers", t, func() {
		var list = []int{1, 2, 3, 4, 5}

		Convey("When I find an element that doesn't exist", func() {
			element, ok := Find(list, func(i int) bool {
				return i == 6
			})

			Convey("Then it should indicate that it was not found", func() {
				So(ok, ShouldBeFalse)
			})

			Convey("Then it should return 0", func() {
				So(element, ShouldEqual, 0)
			})
		})
	})

	Convey("Given an empty list", t, func() {
		var list []int

		Convey("When I find an element", func() {
			element, ok := Find(list, func(i int) bool {
				return i == 3
			})

			Convey("Then it should indicate that it was not found", func() {
				So(ok, ShouldBeFalse)
			})

			Convey("Then it should return 0", func() {
				So(element, ShouldEqual, 0)
			})
		})
	})
}

func TestToString(t *testing.T) {
	Convey("Given a list of integers", t, func() {
		var list = []int{1, 2, 3, 4, 5}

		Convey("When I convert it to a string", func() {
			Convey("Then it should return a string representation if list", func() {
				So(ToString(list), ShouldEqual, "[1 2 3 4 5]")
			})
		})
	})

	Convey("Given an empty list", t, func() {
		var list []int

		Convey("When I convert it to a string", func() {
			Convey("Then it should return an empty string", func() {
				So(ToString(list), ShouldEqual, "[]")
			})
		})
	})

	Convey("Given a type alias to string", t, func() {
		type MyString string
		var myString MyString = "Hello"
		Convey("When I convert it to a string", func() {
			Convey("Then it should return a string value", func() {
				So(ToString(myString), ShouldEqual, "Hello")
			})
		})
	})
}

func TestRemoveIfExists(t *testing.T) {
	Convey("Given a file that exists", t, func() {
		const fileName = "test.txt"
		_, _ = filesystem.Get().Create(fileName)

		Convey("When I remove it", func() {
			err := RemoveIfExists(fileName)

			Convey("Then it should not return an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then it should be removed", func() {
				_, err = os.Stat(fileName)
				So(os.IsNotExist(err), ShouldBeTrue)
			})
		})

		Convey("When I remove it again", func() {
			err := RemoveIfExists(fileName)

			Convey("Then it should not return an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then it should be removed", func() {
				_, err = os.Stat(fileName)
				So(os.IsNotExist(err), ShouldBeTrue)
			})
		})
	})

	Convey("Given a file that doesn't exist", t, func() {
		const fileName = "test.txt"

		Convey("When I remove it", func() {
			err := RemoveIfExists(fileName)

			Convey("Then it should not return an error", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestBytesToMegabytes(t *testing.T) {
	Convey("Given 1024^2 bytes", t, func() {
		const bytes = 1024 * 1024

		Convey("When I convert it to megabytes", func() {
			Convey("Then it should return 1", func() {
				So(math.Floor(BytesToMegabytes(bytes)), ShouldEqual, 1)
			})
		})
	})

	Convey("Given 1025^2 bytes", t, func() {
		const bytes = 1025 * 1025

		Convey("When I convert it to megabytes", func() {
			Convey("Then it should return 1", func() {
				So(math.Floor(BytesToMegabytes(bytes)), ShouldEqual, 1)
			})
		})
	})

	Convey("Given 1023 bytes", t, func() {
		const bytes = 1023

		Convey("When I convert it to megabytes", func() {
			Convey("Then it should return 0", func() {
				So(math.Floor(BytesToMegabytes(bytes)), ShouldEqual, 0)
			})
		})
	})
}

func TestIfElse(t *testing.T) {
	Convey("Given a true", t, func() {
		Convey("Then first value should be returned", func() {
			So(IfElse(true, 1, 2), ShouldEqual, 1)
		})
	})

	Convey("Given a false", t, func() {
		Convey("Then second value should be returned", func() {
			So(IfElse(false, 1, 2), ShouldEqual, 2)
		})
	})
}

func shouldMatch(actual interface{}, expected ...interface{}) string {
	// compile regex
	regex := expected[0].(string)
	r, err := regexp.Compile(regex)

	if err != nil {
		return fmt.Sprintf("Error compiling regex: %s", err)
	}

	// match
	if !r.MatchString(actual.(string)) {
		return fmt.Sprintf("Expected %s to match %s", actual.(string), regex)
	}

	return ""
}

func TestFetchLatestVersion(t *testing.T) {
	Convey("When I fetch the latest version", t, func() {
		version, err := FetchLatestVersion()

		Convey("Then it should not return an error", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then it should return a version", func() {
			So(version, ShouldNotBeEmpty)
		})

		Convey("Then it should look like a version", func() {

			So(version, shouldMatch, `^\d+\.\d+\.\d+$`)
		})
	})
}
