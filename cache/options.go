package cache

import (
	"github.com/samber/mo"
	"time"
)

// Options is a struct that contains options for the cache.
type Options struct {
	// ExpireEvery is the duration after which the cache will be expired.
	// If the value is not specified (mo.None), the cache will never expire.
	ExpireEvery mo.Option[time.Duration]
}
