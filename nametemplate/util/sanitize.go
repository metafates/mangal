//go:build !windows && !darwin

package util

const invalidPathCharsOS = invalidPathCharsUNIX

func sanitizeOS(path string) string {
	return sanitizeUNIX(path)
}
