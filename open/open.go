package open

import (
	"fmt"
	"github.com/metafates/mangal/constant"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	errUnsupportedOS = fmt.Errorf("can't open on this OS: %s", runtime.GOOS)
	runDll32         = filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")
)

func open(input string) (cmd *exec.Cmd, osSupported bool) {
	switch runtime.GOOS {
	case constant.Windows:
		return exec.Command(runDll32, "url.dll,FileProtocolHandler", input), true
	case constant.Darwin:
		return exec.Command("open", input), true
	case constant.Linux:
		return exec.Command("xdg-open", input), true
	case constant.Android:
		return exec.Command("termux-open", input), true
	default:
		return nil, false
	}
}

func openWith(input, with string) (cmd *exec.Cmd, osSupported bool) {
	switch runtime.GOOS {
	case constant.Windows:
		return exec.Command("cmd", "/C", "start", "", with, strings.ReplaceAll(input, "&", "^&")), true
	case constant.Darwin:
		return exec.Command("open", "-a", with, input), true
	case constant.Linux:
		return exec.Command(with, input), true
	case constant.Android:
		return exec.Command("termux-open", "--choose", input), true
	default:
		return nil, false
	}
}
