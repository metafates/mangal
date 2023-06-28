package stringutil

func Trim(s string, max int) string {
	const ellipsis = "…"

	runes := []rune(s)
	if len(runes)-len(ellipsis) >= max {
		return string(runes[:max-len(ellipsis)]) + ellipsis
	}

	return s
}
