package util

import (
	"fmt"
	"github.com/samber/lo"
	"golang.org/x/term"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
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
	{regexp.MustCompile(`[\\/<>:;"'|?!*{}#%&^+,~\s]`), "_"},
	{regexp.MustCompile(`__+`), "_"},
	{regexp.MustCompile(`^[_\-.]+|[_\-.]+$`), ""},
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
	return term.GetSize(int(os.Stdout.Fd()))
}

func FileStem(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

func Wrap(s string, width int) string {
	var lines []string
	for len(s) > width {
		i := strings.LastIndex(s[:width], " ")
		if i == -1 {
			i = width
		}
		lines = append(lines, s[:i])
		s = s[i:]
	}
	lines = append(lines, s)
	return strings.Join(lines, "\n")
}

func ClearScreen() {
	run := func(name string, args ...string) error {
		command := exec.Command(name, args...)
		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		command.Stderr = os.Stderr
		return command.Run()
	}

	switch runtime.GOOS {
	case "linux", "darwin":
		err := run("tput", "clear")
		if err != nil {
			_ = run("clear")
		}
	case "windows":
		_ = run("cls")
	}
}
