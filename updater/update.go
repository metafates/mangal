package updater

import (
	"errors"
	"fmt"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/icon"
	"net/http"
	"os"
	"os/exec"
)

// Update updates mangal to the latest version.
func Update() (err error) {
	method := DetectInstallationMethod()

	switch method {
	case Homebrew:
		fmt.Printf("%s Homebrew installation detected", icon.Get(icon.Progress))
		return updateHomebrew()
	case Scoop:
		fmt.Printf("%s Scoop installation detected", icon.Get(icon.Progress))
		return updateScoop()
	case Termux:
		fmt.Printf("%s Termux installation detected", icon.Get(icon.Progress))
		return updateScript()
	case Script:
		fmt.Printf("%s Script installation detected", icon.Get(icon.Progress))
		return updateScript()
	case Unknown:
		return errors.New("Unknown installation method, can't update")
	}

	return
}

// updateHomebrew updates mangal using Homebrew.
func updateHomebrew() (err error) {
	cmd := exec.Command("brew", "upgrade", constant.Mangal)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// updateScoop updates mangal using Scoop.
func updateScoop() (err error) {
	cmd := exec.Command("Scoop", "update", constant.Mangal)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// updateScript updates mangal using the Script.
func updateScript() (err error) {
	res, err := http.Get(constant.InstallScriptURL)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error fetching Script: status code %d", res.StatusCode)
	}

	var scriptSource []byte
	_, err = res.Body.Read(scriptSource)

	cmd := exec.Command("sh", "-c", string(scriptSource))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
