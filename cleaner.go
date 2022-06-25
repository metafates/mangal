package main

import (
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

	// Check if cache dir exists
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
	tempFiles, err := Afero.ReadDir(tempDir)

	if err != nil {
		return 0, 0
	}

	for _, tempFile := range tempFiles {
		name := tempFile.Name()

		// Check if file is temp file
		if strings.HasPrefix(name, TempPrefix) {

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

	return counter, bytes
}
