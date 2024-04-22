package cache

import (
	"time"
	"sync"
)

type cacheEntry struct {
	createdAt 	time.Time
	val			[]byte
}

type Cache struct {
	mu		sync.Mutex
	cache 	map[string]cacheEntry
}

func newCacheEntry(val []byte) cacheEntry {
	return cacheEntry{time.Now(), val}
}

func (c Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = newCacheEntry(val)
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, exists := c.cache[key]
	return val.val, exists
}