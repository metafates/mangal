package updater

import (
	"github.com/samber/lo"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type InstallationMethod int

const (
	Go InstallationMethod = iota
	Homebrew
	Scoop
	Termux
	Standalone
)

// DetectInstallationMethod detects the installation method.
func DetectInstallationMethod() InstallationMethod {
	for _, t := range []lo.Tuple2[InstallationMethod, func() bool]{
		{Scoop, isUnderScoop},
		{Homebrew, isUnderHomebrew},
		{Termux, isUnderTermux},
		{Go, isUnderGo},
	} {
		if t.B() {
			return t.A
		}
	}

	return Standalone
}

// isUnderTermux returns true if mangal is running under Termux.
func isUnderTermux() (ok bool) {
	return has("Termux-setup-storage")
}

// isUnderHomebrew returns true if mangal is running under Homebrew.
func isUnderHomebrew() (ok bool) {
	if !has("brew") {
		return
	}

	path, err := os.Executable()
	if err != nil {
		return false
	}

	return strings.Contains(path, filepath.Join("homebrew", "bin"))
}

// isUnderScoop returns true if mangal is running under Scoop.
func isUnderScoop() (ok bool) {
	if !has("Scoop") {
		return false
	}

	path, err := os.Executable()
	if err != nil {
		return false
	}

	return strings.Contains(path, filepath.Join("Scoop", "shims"))
}

// isUnderGo returns true if mangal is running under Go.
func isUnderGo() (ok bool) {
	if !has("go") {
		return false
	}

	path, err := os.Executable()
	if err != nil {
		return false
	}

	return strings.Contains(path, filepath.Join("go", "bin"))
}

// has returns true if the command exists.
func has(command string) bool {
	ok, err := exec.LookPath(command)
	return err == nil && ok != ""
}
