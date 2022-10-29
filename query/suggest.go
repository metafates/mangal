package query

import (
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/samber/lo"
	"github.com/samber/mo"
)

var (
	suggestionCache = make(map[string]mo.Option[string])
)

// Suggest gives a suggestion for a query
func Suggest(query string) mo.Option[string] {
	query = sanitize(query)

	if prev, ok := suggestionCache[query]; ok {
		return prev
	}

	cached, ok := cacher.Get().Get()
	if !ok {
		return mo.None[string]()
	}

	// fuzzy filter keys
	// and get the one with the highest rank
	matching := fuzzy.Find(query, lo.Keys(cached))

	result := lo.MaxBy(matching, func(a, b string) bool {
		return cached[a].Rank > cached[b].Rank
	})

	var suggestion mo.Option[string]

	if result == "" {
		suggestion = mo.None[string]()
	} else {
		suggestion = mo.Some(result)
	}

	suggestionCache[query] = suggestion

	return suggestion
}
