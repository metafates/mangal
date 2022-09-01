package updater

import (
	"errors"
	"fmt"
	"github.com/metafates/mangal/icon"
	"net/http"
	"os"
	"os/exec"
)

// Update updates mangal to the latest version.
func Update() (err error) {
	method := detectInstallationMethod()

	switch method {
	case homebrew:
		fmt.Printf("%s Homebrew installation detected", icon.Get(icon.Progress))
		return updateHomebrew()
	case scoop:
		fmt.Printf("%s Scoop installation detected", icon.Get(icon.Progress))
		return updateScoop()
	case termux:
		fmt.Printf("%s Termux installation detected", icon.Get(icon.Progress))
		return updateScript()
	case script:
		fmt.Printf("%s Script installation detected", icon.Get(icon.Progress))
		return updateScript()
	case unknown:
		return errors.New("unknown installation method, can't update")
	}

	return
}

// updateHomebrew updates mangal using homebrew.
func updateHomebrew() (err error) {
	cmd := exec.Command("brew", "upgrade", "mangal")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// updateScoop updates mangal using scoop.
func updateScoop() (err error) {
	cmd := exec.Command("scoop", "update", "mangal")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// updateScript updates mangal using the script.
func updateScript() (err error) {
	res, err := http.Get("https://raw.githubusercontent.com/metafates/mangal/main/scripts/install")
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error fetching script: status code %d", res.StatusCode)
	}

	var scriptSource []byte
	_, err = res.Body.Read(scriptSource)

	cmd := exec.Command("sh", "-c", string(scriptSource))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
