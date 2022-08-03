package util

import (
	"github.com/samber/lo"
	"regexp"
)

func PrettyTrim(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return s[:maxLen] + "..."
}

func PadZero(s string, l int) string {
	for l > len(s) {
		s = "0" + s
	}

	return s
}

// replacers is a list of regexp.Regexp pairs that will be used to sanitize filenames.
var replacers = []lo.Tuple2[*regexp.Regexp, string]{
	{regexp.MustCompile(`[\\/<>:"|?*\s]`), "_"},
	{regexp.MustCompile(`__+`), "_"},
	{regexp.MustCompile(`^_+|_+$`), ""},
	{regexp.MustCompile(`^\.+|\.+$`), ""},
}

// SanitizeFilename will remove all invalid characters from a path.
func SanitizeFilename(filename string) string {
	for _, re := range replacers {
		filename = re.A.ReplaceAllString(filename, re.B)
	}

	return filename
}
