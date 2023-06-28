package path

import (
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/mangalorg/mangal/meta"
)

func HomeDir() string {
    dir := xdg.Home
    createIfAbsent(dir)
    return dir
}

func CacheDir() string {
    dir := filepath.Join(xdg.CacheHome, meta.AppName)
    createIfAbsent(dir)
    return dir
}

func ConfigDir() string {
    dir := filepath.Join(xdg.ConfigHome, meta.AppName)
    createIfAbsent(dir)
    return dir
}

func DownloadsDir() string {
    dir := xdg.UserDirs.Download 
    createIfAbsent(dir)
    return dir
}
