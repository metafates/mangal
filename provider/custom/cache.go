package custom

import (
	"github.com/metafates/mangal/cache"
	"github.com/samber/mo"
	"time"
)

type cacher[T any] struct {
	internal *cache.Cache[map[string]T]
}

func newCacher[T any](name string) *cacher[T] {
	return &cacher[T]{
		internal: cache.New[map[string]T](name, &cache.Options{
			ExpireEvery: mo.Some(time.Hour * 24),
		}),
	}
}

func (c *cacher[T]) Get(key string) mo.Option[T] {
	if c.internal.Get().IsAbsent() {
		return mo.None[T]()
	}

	data := c.internal.Get().MustGet()

	if x, ok := data[key]; ok {
		return mo.Some(x)
	}

	return mo.None[T]()
}

func (c *cacher[T]) Set(key string, t T) error {
	var data map[string]T

	if c.internal.Get().IsAbsent() {
		data = make(map[string]T)
	} else {
		data = c.internal.Get().MustGet()
	}

	data[key] = t

	return c.internal.Set(data)
}
