package updater

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func Update() (err error) {
	method := detectInstallationMethod()

	switch method {
	case homebrew:
		return updateHomebrew()
	case scoop:
		return updateScoop()
	case script:
		return updateScript()
	case unknown:
		return errors.New("unknown installation method, can't update")
	}

	return
}

func updateHomebrew() (err error) {
	cmd := exec.Command("brew", "upgrade", "mangal")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func updateScoop() (err error) {
	cmd := exec.Command("scoop", "update", "mangal")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func updateScript() (err error) {
	res, err := http.Get("https://raw.githubusercontent.com/metafates/mangal/main/install")
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
