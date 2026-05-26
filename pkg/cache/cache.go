package cache

import (
	"sync"
	"time"
)

type item struct {
	value      interface{}
	expiration int64
}

type Cache struct {
	items map[string]item
	mu    sync.RWMutex
}

func NewCache() *Cache {
	c := &Cache{
		items: make(map[string]item),
	}
	
	// Start a cleanup goroutine to remove expired items periodically
	go c.cleanupLoop()
	return c
}

func (c *Cache) Set(key string, value interface{}, d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var exp int64
	if d > 0 {
		exp = time.Now().Add(d).UnixNano()
	}

	c.items[key] = item{
		value:      value,
		expiration: exp,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
		return nil, false
	}

	return item.value, true
}

func (c *Cache) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		c.cleanup()
	}
}

func (c *Cache) cleanup() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if v.expiration > 0 && now > v.expiration {
			delete(c.items, k)
		}
	}
}
