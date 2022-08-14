package where

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/samber/lo"
	"os"
	"path/filepath"
)

const EnvConfigPath = "MANGAL_CONFIG_PATH"

func mkdir(path string) string {
	lo.Must0(filesystem.Get().MkdirAll(path, os.ModePerm))
	return path
}

func Config() string {
	var path string

	customDir, present := os.LookupEnv(EnvConfigPath)
	if present {
		path = customDir
	} else {
		path = filepath.Join(lo.Must(os.UserConfigDir()), constant.Mangal)
	}

	return mkdir(path)
}

func Sources() string {
	return mkdir(filepath.Join(Config(), "sources"))
}

func Logs() string {
	return mkdir(filepath.Join(Config(), "logs"))
}

func History() string {
	cacheDir := filepath.Join(lo.Must(os.UserCacheDir()), constant.CachePrefix)
	lo.Must0(filesystem.Get().MkdirAll(filepath.Dir(cacheDir), os.ModePerm))

	path := filepath.Join(mkdir(cacheDir), "history.json")

	exists := lo.Must(filesystem.Get().Exists(path))
	if !exists {
		lo.Must0(filesystem.Get().WriteFile(path, []byte("{}"), os.ModePerm))
	}

	return path
}
