package util

import (
	"strings"
	"unicode"
)

const (
	invalidPathCharsUNIX    = `/`
	invalidPathCharsDarwin  = `/:`
	invalidPathCharsWindows = `<>:"/\|?*`
)

func sanitizePath(path string, isInvalid func(rune) bool) string {
	var (
		sanitized strings.Builder
		prev      rune
	)

	const underscore = '_'

	for _, r := range path {
		var toWrite rune
		if isInvalid(r) {
			toWrite = underscore
		} else {
			toWrite = r
		}

		// replace two or more consecutive underscores with one underscore
		if (toWrite == underscore && prev != underscore) || toWrite != underscore {
			sanitized.WriteRune(toWrite)
		}

		prev = toWrite
	}

	return strings.TrimFunc(sanitized.String(), func(r rune) bool {
		return r == underscore || unicode.IsSpace(r)
	})
}

func sanitizeWhitespace(path string) string {
	return sanitizePath(path, func(r rune) bool {
		return unicode.IsSpace(r)
	})
}

func sanitizeUNIX(path string) string {
	return sanitizePath(path, func(r rune) bool {
		return strings.ContainsRune(invalidPathCharsUNIX, r)
	})
}

func sanitizeDarwin(path string) string {
	return sanitizePath(path, func(r rune) bool {
		return strings.ContainsRune(invalidPathCharsDarwin, r)
	})
}

func sanitizeWindows(path string) string {
	return sanitizePath(path, func(r rune) bool {
		return strings.ContainsRune(invalidPathCharsWindows, r)
	})
}
