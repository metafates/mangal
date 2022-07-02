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

	files, err := Afero.ReadDir(cacheDir)

	for _, cacheFile := range files {
		name := cacheFile.Name()

		// Check if file is cache file
		if strings.HasPrefix(name, CachePrefix) {

			p := filepath.Join(cacheDir, name)

			if cacheFile.IsDir() {
				b, err := DirSize(p)
				if err == nil {
					bytes += b
				}
			}

			err = Afero.RemoveAll(p)
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
