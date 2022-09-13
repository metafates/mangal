package where

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const EnvConfigPath = "MANGAL_CONFIG_PATH"

// mkdir creates a directory and all parent directories if they don't exist
// will return the path of the directory
func mkdir(path string) string {
	lo.Must0(filesystem.Api().MkdirAll(path, os.ModePerm))
	return path
}

// Config path
// Will create the directory if it doesn't exist
func Config() string {
	var path string

	if customDir, present := os.LookupEnv(EnvConfigPath); present {
		path = customDir
	} else {
		path = filepath.Join(lo.Must(os.UserConfigDir()), constant.Mangal)
	}

	return mkdir(path)
}

// Sources path
// Will create the directory if it doesn't exist
func Sources() string {
	return mkdir(filepath.Join(Config(), "sources"))
}

// Logs path
// Will create the directory if it doesn't exist
func Logs() string {
	return mkdir(filepath.Join(Config(), "logs"))
}

// History path to the file
// Will create the directory if it doesn't exist
func History() string {
	genericCacheDir, err := os.UserCacheDir()
	if err != nil {
		genericCacheDir = "."
	}

	path := filepath.Join(genericCacheDir, constant.CachePrefix+"history.json")

	exists := lo.Must(filesystem.Api().Exists(path))
	if !exists {
		lo.Must0(filesystem.Api().WriteFile(path, []byte("{}"), os.ModePerm))
	}

	return path
}

// Downloads path
// Will create the directory if it doesn't exist
func Downloads() string {
	path, err := filepath.Abs(viper.GetString(constant.DownloaderPath))

	if err != nil {
		path, err = os.Getwd()
		if err != nil {
			path = "."
		}
	}

	return mkdir(path)
}

// Cache path
// Will create the directory if it doesn't exist
func Cache() string {
	genericCacheDir, err := os.UserCacheDir()
	if err != nil {
		genericCacheDir = "."
	}

	cacheDir := filepath.Join(genericCacheDir, constant.CachePrefix)
	return mkdir(cacheDir)
}

// Temp path
// Will create the directory if it doesn't exist
func Temp() string {
	tempDir := filepath.Join(os.TempDir(), constant.TempPrefix)
	return mkdir(tempDir)
}
