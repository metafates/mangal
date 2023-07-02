package stringutil

import "fmt"

func Trim(s string, max int) string {
	const ellipsis = "…"

	runes := []rune(s)
	if len(runes)-len(ellipsis) >= max {
		return string(runes[:max-len(ellipsis)]) + ellipsis
	}

	return s
}

func Quantify(n int, singular, plural string) string {
	var form string
	if n == 1 {
		form = singular
	} else {
		form = plural
	}

	return fmt.Sprint(n, " ", form)
}
