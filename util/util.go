package util

import (
	"fmt"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
	"golang.org/x/term"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// PadZero pads a number with leading zeros.
func PadZero(s string, l int) string {
	return strings.Repeat("0", Max(l-len(s), 0)) + s
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

// Quantify returns a string with the given number and unit.
func Quantify(count int, singular, plural string) string {
	if count == 1 {
		return fmt.Sprintf("%d %s", count, singular)
	}

	return fmt.Sprintf("%d %s", count, plural)
}

// TerminalSize returns the dimensions of the given terminal.
func TerminalSize() (width, height int, err error) {
	return term.GetSize(int(os.Stdout.Fd()))
}

// FileStem returns the file name without the extension.
func FileStem(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

// ClearScreen clears the terminal screen.
func ClearScreen() {
	run := func(name string, args ...string) error {
		command := exec.Command(name, args...)
		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		command.Stderr = os.Stderr
		return command.Run()
	}

	switch runtime.GOOS {
	case constant.Linux, constant.Darwin:
		err := run("tput", "clear")
		if err != nil {
			_ = run("clear")
		}
	case constant.Windows:
		_ = run("cls")
	}
}

// ReGroups parses the string with the given regular expression and returns the
// group values defined in the expression.
func ReGroups(pattern *regexp.Regexp, str string) (groups map[string]string) {
	groups = make(map[string]string)
	match := pattern.FindStringSubmatch(str)

	for i, name := range pattern.SubexpNames() {
		if i > 0 && i <= len(match) {
			groups[name] = match[i]
		}
	}

	return
}

// Ignore calls function and explicitely ignores error
func Ignore(f func() error) {
	_ = f()
}

// Max returns the maximum value of the given items.
func Max[T constraints.Ordered](items ...T) (max T) {
	for _, item := range items {
		if item > max {
			max = item
		}
	}

	return
}

// Min returns the minimum value of the given items.
func Min[T constraints.Ordered](items ...T) (min T) {
	min = items[0]
	for _, item := range items {
		if item < min {
			min = item
		}
	}

	return
}

// PrintErasable prints a string that can be erased by calling a returned function.
func PrintErasable(msg string) (eraser func()) {
	_, _ = fmt.Fprintf(os.Stdout, "\r%s", msg)

	return func() {
		_, _ = fmt.Fprintf(os.Stdout, "\r%s\r", strings.Repeat(" ", len(msg)))
	}
}

// Capitalize returns a string with the first letter capitalized.
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

// Delete removes the given path from the filesystem.
// It can handle both files and directories (recursively).
func Delete(path string) error {
	stat, err := filesystem.Api().Stat(path)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return filesystem.Api().RemoveAll(path)
	}

	return filesystem.Api().Remove(path)
}
