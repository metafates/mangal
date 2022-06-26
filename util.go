package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/constraints"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// IfElse is a ternary operator equavlient
// 	Example: IfElse(true, "true", "false") => "true"
func IfElse[T any](condition bool, then, othwerwise T) T {
	if condition {
		return then
	}

	return othwerwise
}

// Contains checks if slice contains element
// 	Example: Contains([]int{1, 2, 3}, 2) => true
func Contains[T comparable](slice []T, elem T) bool {
	for _, el := range slice {
		if el == elem {
			return true
		}
	}

	return false
}

// BytesToMegabytes converts bytes to megabytes
// 	Example: BytesToMegabytes(1024 * 1024) => 1
func BytesToMegabytes(bytes int64) float64 {
	return float64(bytes) / 1_048_576 // 1024 ^ 2
}

// PrettyTrim trims string to given size and adds ellipsis to the end
// 	Example: PrettyTrim("Hello World", 5) => "Hello..."
func PrettyTrim(text string, limit int) string {
	if len(text)-3 > limit {
		return text[:limit-3] + "..."
	}

	return text
}

// Plural makes singular word a plural if count ≠ 1
// 	Example: Plural("book", 2) => "books"
func Plural(word string, count int) string {
	if count == 1 {
		return word
	}

	return word + "s"
}

// Max between 2 values
// 	Example: Max(1, 2) => 2
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}

	return b
}

// IsUnique tests if given array has only unique elements (e.g. if it's a set)
// 	Example: IsUnique([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) => true
func IsUnique[T comparable](arr []T) bool {
	for i, a := range arr {
		for j, b := range arr {
			if i == j {
				continue
			}

			if a == b {
				return false
			}
		}
	}
	return true
}

// DirSize gets directory size in bytes
// 	Example: DirSize("/tmp") => 1234
func DirSize(path string) (int64, error) {
	var size int64
	err := Afero.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

// Find element in slice by function
// 	Example: Find([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, func(i int) bool { return i > 5 }) => 6
func Find[T any](list []T, f func(T) bool) (T, bool) {
	var prev *T

	for _, el := range list {
		prev = &el
		if f(el) {
			return el, true
		}
	}

	return *prev, false
}

// Map applies function to each element of the slice
// 	Example: Map([]int{1, 2, 3}, func(i int) int { return i * 2 }) => [2, 4, 6]
func Map[T, G any](list []T, f func(T) G) []G {
	var mapped = make([]G, len(list))

	for i, el := range list {
		mapped[i] = f(el)
	}

	return mapped
}

// ToString converts any type to string
// 	Example: ToString(1) => "1"
func ToString[T any](v T) string {
	return fmt.Sprintf("%v", v)
}

// FetchLatestVersion fetches the latest version tag of the app from the GitHub releases
// 	Example: FetchLatestVersion() => "1.0.0"
func FetchLatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/metafates/mangal/releases/latest")
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	// remove the v from the tag name
	return release.TagName[1:], nil
}

// SanitizePath will remove all invalid characters from a path
func SanitizePath(path string) string {
	const forbiddenChars = `\\/<>:"|?*`

	// remove all forbidden characters using regex
	path = regexp.MustCompile(`[`+forbiddenChars+`]`).ReplaceAllString(path, "")

	// replace all spaces with underscores
	path = strings.ReplaceAll(path, " ", "_")

	// remove all double underscores
	path = regexp.MustCompile(`__+`).ReplaceAllString(path, "_")

	// remove all leading and trailing underscores
	path = regexp.MustCompile(`^_+|_+$`).ReplaceAllString(path, "")

	// remove all leading and trailing dots
	path = regexp.MustCompile(`^\.+|\.+$`).ReplaceAllString(path, "")

	return path
}
