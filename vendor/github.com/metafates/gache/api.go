package gache

import (
	"time"
)

// Set sets the value of the cache.
// If initialization or marshalling fails, it will return an error.
// In memory-only mode it will never fail.
// It will restart the cache's lifetime.
func (g *Cache[T]) Set(value T) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	err := g.init()
	if err != nil {
		return err
	}

	// update value
	g.data.Internal = value

	// update time
	now := time.Now()
	g.data.Time = &now

	if g.options.Path != "" {
		err = g.save()
	}

	return err
}

// Get returns the value of the cache.
// If initialization fails, it will return an error.
// In memory-only mode it will never fail.
func (g *Cache[T]) Get() (cached T, expired bool, err error) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	err = g.init()
	if err != nil {
		return
	}

	// Do not use tryExpire() here, because it modifies the cache, but we RLocked mutex.
	if g.isExpired() {
		expired = true
		return
	}

	cached = g.data.Internal
	return
}
