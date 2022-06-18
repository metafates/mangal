package main

import (
	"os"
	"path/filepath"
	"strings"
)

func RemoveCache() (int, int64) {
	var (
		// counter of removed files
		counter int
		// bytes removed
		bytes int64
	)

	// Cleanup cache files
	cacheDir, err := os.UserCacheDir()
	if err == nil {
		scraperCacheDir := filepath.Join(cacheDir, CachePrefix)
		if exists, err := Afero.Exists(scraperCacheDir); err == nil && exists {
			files, err := Afero.ReadDir(scraperCacheDir)
			if err == nil {
				for _, f := range files {
					counter++
					bytes += f.Size()
				}
			}

			_ = Afero.RemoveAll(scraperCacheDir)
		}
	}

	return counter, bytes
}

func RemoveTemp() (int, int64) {
	var (
		// counter of removed files
		counter int
		// bytes removed
		bytes int64
	)

	// Cleanup temp files
	tempDir := os.TempDir()
	tempFiles, err := Afero.ReadDir(tempDir)
	if err == nil {
		lowerAppName := strings.ToLower(AppName)
		for _, tempFile := range tempFiles {
			name := tempFile.Name()
			if strings.HasPrefix(name, AppName) || strings.HasPrefix(name, lowerAppName) {

				p := filepath.Join(tempDir, name)

				if tempFile.IsDir() {
					b, err := DirSize(p)
					if err == nil {
						bytes += b
					}
				}

				err = Afero.RemoveAll(p)
				if err == nil {
					bytes += tempFile.Size()
					counter++
				}
			}
		}
	}

	return counter, bytes
}
