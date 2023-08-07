package util

const invalidPathCharsOS = invalidPathCharsWindows

func sanitizeOS(path string) string {
	return sanitizeWindows(path)
}
