package open

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	unsupportedOSError = fmt.Errorf("can't open on this OS: %s", runtime.GOOS)
	runDll32           = filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")
	termuxOpen         = "env -i PATH=/data/data/com.termux/files/usr/bin ANDROID_DATA=/data am broadcast -a android.intent.action.VIEW -n com.termux/com.termux.app.TermuxOpenReceiver -d"
)

func open(input string) (cmd *exec.Cmd, osSupported bool) {
	switch runtime.GOOS {
	case "windows":
		return exec.Command(runDll32, "url.dll,FileProtocolHandler", input), true
	case "darwin":
		return exec.Command("open", input), true
	case "linux":
		return exec.Command("xdg-open", input), true
	case "android":
		return exec.Command(termuxOpen, input), true
	default:
		return nil, false
	}
}

func openWith(input, with string) (cmd *exec.Cmd, osSupported bool) {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("cmd", "/C", "start", "", with, strings.ReplaceAll(input, "&", "^&")), true
	case "darwin":
		return exec.Command("open", "-a", with, input), true
	case "linux":
		return exec.Command(with, input), true
	case "android":
		return exec.Command(termuxOpen, input), true
	default:
		return nil, false
	}
}
