package mini

import "github.com/samber/lo"

type bind lo.Tuple2[string, string]

var (
	quit   = &bind{A: "q", B: "quit"}
	prev   = &bind{A: "p", B: "previous"}
	next   = &bind{A: "n", B: "next"}
	reread = &bind{A: "r", B: "reread"}
	back   = &bind{A: "b", B: "back"}
	search = &bind{A: "s", B: "search"}
)

func (b *bind) eq(other *bind) bool {
	return other != nil && b.A == other.A
}

func parseBind(b string) (*bind, bool) {
	switch b {
	case quit.A:
		return quit, true
	case prev.A:
		return prev, true
	case next.A:
		return next, true
	case reread.A:
		return reread, true
	case back.A:
		return back, true
	case search.A:
		return search, true
	default:
		return nil, false
	}
}
