package cleaner

import (
	"github.com/metafates/mangal/common"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/util"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"strings"
)

// RemoveCache removes cache files
func RemoveCache() (int, int64) {
	var (
		// counter of removed files
		counter int
		// bytes removed
		bytes int64
	)

	// Cleanup cache files
	cacheDir, err := os.UserCacheDir()

	if err != nil {
		return 0, 0
	}

	files, err := afero.ReadDir(filesystem.Get(), cacheDir)

	for _, cacheFile := range files {
		name := cacheFile.Name()

		// Check if file is cache file
		if strings.HasPrefix(name, common.CachePrefix) {

			p := filepath.Join(cacheDir, name)

			if cacheFile.IsDir() {
				b, err := util.DirSize(p)
				if err == nil {
					bytes += b
				}
			}

			err = filesystem.Get().RemoveAll(p)
			if err == nil {
				bytes += cacheFile.Size()
				counter++
			}
		}
	}

	return counter, bytes
}

// RemoveTemp removes temp files
func RemoveTemp() (int, int64) {
	var (
		// counter of removed files
		counter int
		// bytes removed
		bytes int64
	)

	tempDir := os.TempDir()
	tempFiles, err := afero.ReadDir(filesystem.Get(), tempDir)

	if err != nil {
		return 0, 0
	}

	for _, tempFile := range tempFiles {
		name := tempFile.Name()

		// Check if file is temp file
		if strings.HasPrefix(name, common.TempPrefix) {

			p := filepath.Join(tempDir, name)

			if tempFile.IsDir() {
				b, err := util.DirSize(p)
				if err == nil {
					bytes += b
				}
			}

			err = filesystem.Get().RemoveAll(p)
			if err == nil {
				bytes += tempFile.Size()
				counter++
			}
		}
	}

	return counter, bytes
}

// RemoveHistory removes history file
func RemoveHistory() (int, int64) {
	path, err := util.HistoryFilePath()

	if err != nil {
		return 0, 0
	}

	// get file size
	stat, err := filesystem.Get().Stat(path)

	if err != nil {
		return 0, 0
	}

	err = util.RemoveIfExists(path)

	if err != nil {
		return 0, 0
	}

	return 1, stat.Size()
}
