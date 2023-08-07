package util

const invalidPathCharsOS = invalidPathCharsDarwin

func sanitizeOS(path string) string {
	return sanitizeDarwin(path)
}
