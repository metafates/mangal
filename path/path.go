package path

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"
	"github.com/mangalorg/mangal/meta"
)

func HomeDir() string {
	dir := xdg.Home
	createDirIfAbsent(dir)
	return dir
}

func CacheDir() string {
	dir := filepath.Join(xdg.CacheHome, meta.AppName)
	createDirIfAbsent(dir)
	return dir
}

func ConfigDir() string {
	var dir string

	if runtime.GOOS == "darwin" {
		dir = filepath.Join(xdg.Home, ".config", meta.AppName)
	} else {
		dir = filepath.Join(xdg.ConfigHome, meta.AppName)
	}

	createDirIfAbsent(dir)
	return dir
}

func DownloadsDir() string {
	dir := xdg.UserDirs.Download
	createDirIfAbsent(dir)
	return dir
}

func TempDir() string {
	dir := filepath.Join(os.TempDir(), meta.AppName)
	createDirIfAbsent(dir)
	return dir
}

func ProvidersDir() string {
	dir := filepath.Join(ConfigDir(), "providers")
	createDirIfAbsent(dir)
	return dir
}

func LuaProvidersDir() string {
	dir := filepath.Join(ProvidersDir(), "lua")
	createDirIfAbsent(dir)
	return dir
}
