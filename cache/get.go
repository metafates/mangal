package cache

import "github.com/samber/mo"

// Get returns the cached data if it exists and is not expired, otherwise none.
func (c *Cache[T]) Get() mo.Option[T] {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	_ = c.init()

	return c.data.Internal
}
