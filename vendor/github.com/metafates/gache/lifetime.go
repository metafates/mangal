package gache

import "time"

func (g *Cache[T]) isExpired() bool {
	return g.options.Lifetime >= 0 &&
		g.data.Time != nil &&
		time.Since(*g.data.Time) > g.options.Lifetime
}

func (g *Cache[T]) tryExpire() error {
	// check if the cache has expired
	if g.isExpired() {
		// erase the cache
		var defaultT T
		g.data = &chronoData[T]{
			Internal: defaultT,
			Time:     nil,
		}

		// call the expiration hook
		if g.options.ExpirationHook != nil {
			g.options.ExpirationHook()
		}

		// save the cache
		return g.save()
	}

	return nil
}
