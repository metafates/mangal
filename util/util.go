package util

import (
	"fmt"
	"github.com/samber/lo"
	"golang.org/x/sys/unix"
	"regexp"
	"strings"
)

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

func Quantity(count int, thing string) string {
	if strings.HasSuffix(thing, "s") {
		thing = thing[:len(thing)-1]
	}

	if count == 1 {
		return fmt.Sprintf("%d %s", count, thing)
	}

	return fmt.Sprintf("%d %ss", count, thing)
}

// TerminalSize returns the dimensions of the given terminal.
func TerminalSize() (width, height int, err error) {
	fd := unix.Stdout
	ws, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	if err != nil {
		return -1, -1, err
	}
	return int(ws.Col), int(ws.Row), nil
}

// Capitalize will capitalize the first letter of a string.
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}
