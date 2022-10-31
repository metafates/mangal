package mangadex

import (
	"github.com/metafates/mangal/cache"
	"github.com/metafates/mangal/where"
	"github.com/samber/mo"
	"path/filepath"
	"time"
)

type cacher[T any] struct {
	internal *cache.Cache[map[string]T]
}

func newCacher[T any](name string) *cacher[T] {
	return &cacher[T]{
		internal: cache.New[map[string]T](
			filepath.Join(where.Cache(), name+".json"),
			&cache.Options{
				ExpireEvery: mo.Some(time.Hour * 24),
			},
		),
	}
}

func (c *cacher[T]) Get(key string) mo.Option[T] {
	cached, ok := c.internal.Get().Get()
	if !ok {
		return mo.None[T]()
	}

	if value, ok := cached[key]; ok {
		return mo.Some[T](value)
	}

	return mo.None[T]()
}

func (c *cacher[T]) Set(key string, value T) error {
	cached, ok := c.internal.Get().Get()
	if !ok {
		cached = map[string]T{}
	}

	cached[key] = value
	return c.internal.Set(cached)
}
