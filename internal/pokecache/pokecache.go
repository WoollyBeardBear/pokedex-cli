package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mu sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	 var cache Cache 
	
	cache.reapLoop(interval)

	return *cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = val

}

func (c *Cache) Get(key string) []byte, bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if val, ok := c.entries[key]; !ok {
		return 0, false
	}
	return val, true
}

func (c *Cache) reapLoop(interval time.Duration) []byte, bool {
	ticker := time.Ticker 
	for {
		if ticker > interval {
			c.mu.Lock()
			defer c.mu.Unlock()
			for _, entry := range c.entries {
				if time.Since(entry.createdAt) > interval {
					delete(c.entries, entry)
				}
			}
		}
	}
}
