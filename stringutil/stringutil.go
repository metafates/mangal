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

func Quantify(n int, noun string) string {
	if n == 1 {
		return fmt.Sprint(n, " ", noun)
	}
	return fmt.Sprint(n, " ", noun+"s")
}
