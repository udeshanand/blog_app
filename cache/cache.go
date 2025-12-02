package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type Cache struct {
	Store *cache.Cache
}

var DefaultExpiration = 5 * time.Minute

func New() *Cache {
	return &Cache{
		Store: cache.New(DefaultExpiration, 10*time.Minute),
	}
}

func (c *Cache) Set(key string, value interface{}, d time.Duration) {
	c.Store.Set(key, value, d)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	return c.Store.Get(key)
}

func (c *Cache) Delete(key string) {
	c.Store.Delete(key)
}

func (c *Cache) Flush() {
	c.Store.Flush()
}
