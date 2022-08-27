package updater

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func Update() (err error) {
	method := detectInstallationMethod()

	switch method {
	case golang:
		return updateGo()
	case homebrew:
		fmt.Println("homebrew")
	case scoop:
		fmt.Println("scoop")
	case script:
		fmt.Println("script")
	case unknown:
		return errors.New("unknown installation method, can't update")
	}

	return
}

func updateGo() (err error) {
	cmd := exec.Command("go", "install", "github.com/metafates/mangal@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
