package updater

import (
	"errors"
	"fmt"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/util"
	"github.com/metafates/mangal/where"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var erase func()

func info(format string, args ...any) (erase func()) {
	return util.PrintErasable(
		fmt.Sprintf(
			"%s %s",
			icon.Get(icon.Progress),
			fmt.Sprintf(format, args...),
		),
	)
}

// Update updates mangal to the latest version.
func Update() (err error) {
	erase = info("Fetching latest version")
	version, err := LatestVersion()
	if err != nil {
		return
	}

	erase()
	if constant.Version >= version {
		fmt.Printf(
			"%s %s %s\n",
			style.Green("Congrats!"),
			"You're already on the latest version of "+constant.Mangal,
			style.Faint("(which is "+constant.Version+")"),
		)
		return
	}

	method := DetectInstallationMethod()

	switch method {
	case Homebrew:
		erase = info("Homebrew installation detected")
		err = updateHomebrew()
	case Scoop:
		erase = info("Scoop installation detected")
		err = updateScoop()
	case Termux, Go, Standalone:
		erase = info("Non-package manager installation detected")
		err = update()
	default:
		err = errors.New("unknown installation method, can't update")
		return
	}

	if err != nil {
		return
	}

	fmt.Printf(`Updated.

%s

Report any bugs:

    %s

What's new:

    %s

Changelog:

    %s
`,
		style.Combined(style.Bold, style.Green)("Welcome to mangal v"+version),
		style.Faint("https://github.com/metafates/mangal/issues"),
		style.Cyan("https://github.com/metafates/mangal/releases/tag/v"+version),
		style.Faint(fmt.Sprintf("https://github.com/metafates/mangal/compare/v%s...v%s", constant.Version, version)),
	)

	return
}

// updateHomebrew updates mangal using Homebrew.
func updateHomebrew() (err error) {
	erase()
	info("Running %s", style.Yellow("brew upgrade mangal"))
	cmd := exec.Command("brew", "upgrade", constant.Mangal)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// updateScoop updates mangal using Scoop.
func updateScoop() (err error) {
	erase()
	info("Running %s", style.Yellow("scoop update mangal"))
	cmd := exec.Command("Scoop", "update", constant.Mangal)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// update mangal by downloading it directly.
func update() (err error) {
	erase()
	log.Info("Updating mangal to the latest version")

	var (
		version     string
		arch        string
		selfPath    string
		archiveName string
		archiveType string
	)

	if selfPath, err = os.Executable(); err != nil {
		return
	}

	if version, err = LatestVersion(); err != nil {
		return
	}

	switch runtime.GOARCH {
	case "amd64":
		arch = "x86_64"
	case "386":
		arch = "i386"
	default:
		arch = runtime.GOARCH
	}

	if runtime.GOOS == "windows" {
		archiveType = "zip"
	} else {
		archiveType = "tar.gz"
	}

	archiveName = fmt.Sprintf("%s_%s_%s_%s.%s", constant.Mangal, version, util.Capitalize(runtime.GOOS), arch, archiveType)
	url := fmt.Sprintf(
		"https://github.com/metafates/%s/releases/download/v%s/%s",
		constant.Mangal,
		version,
		archiveName,
	)

	erase = info("Downloading %s", style.Yellow(url))

	res, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return err
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("error downloading binary: status code %d", res.StatusCode)
		log.Error(err)
		return
	}

	defer util.Ignore(res.Body.Close)

	erase()
	erase = info("Extracting %s", style.Yellow(archiveName))
	out := filepath.Join(where.Temp(), "mangal_update")

	switch archiveType {
	case "zip":
		archive, err := io.ReadAll(res.Body)
		if err != nil {
			log.Error(err)
			return err
		}

		archivePath := filepath.Join(out, archiveName)
		file, err := filesystem.Api().OpenFile(archivePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		if err != nil {
			log.Error(err)
			return err
		}

		defer util.Ignore(file.Close)

		_, err = file.Write(archive)
		if err != nil {
			log.Error(err)
			return err
		}

		stat, err := file.Stat()
		if err != nil {
			log.Error(err)
			return err
		}

		err = util.Unzip(file, stat.Size(), out)
	case "tar.gz":
		err = util.UntarGZ(res.Body, out)
	}

	if err != nil {
		log.Error(err)
		return err
	}

	erase()
	erase = info("Replacing %s", style.Yellow(selfPath))
	// remove the old binary
	// it should not interrupt the running process
	file, err := filesystem.Api().OpenFile(selfPath, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Error(err)
		err = errors.New("error removing old binary, try running this as a root user")
		return err
	}

	newMangal, err := filesystem.Api().OpenFile(filepath.Join(out, constant.Mangal), os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error(err)
		return
	}

	stat, err := newMangal.Stat()
	if err != nil {
		log.Error(err)
		return
	}

	if stat.Size() == 0 {
		log.Error(err)
		return err
	}

	_, err = io.Copy(file, newMangal)
	if err != nil {
		log.Error(err)
		return err
	}

	erase()

	return
}
