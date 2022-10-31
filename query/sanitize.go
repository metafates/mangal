package query

import "strings"

func sanitize(query string) string {
	return strings.TrimSpace(strings.ToLower(query))
}
