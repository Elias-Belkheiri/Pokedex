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

func (c Cache) reapLoop(interval time.Duration) {
	timeNow := time.Now()
	for k, v := range(c.cache) {
		if timeNow.Sub(v.createdAt) > interval {
			delete(c.cache, k)
		}
	}
}

// func main() {
// 	cache := Cache{cache : make(map[string]cacheEntry)}

// 	val :=  []byte{1,2,3}
// 	cache.Add("hola", val)

// 	val, ok := cache.Get("hola")
// 	fmt.Println(val, ok)

// 	val, ok = cache.Get("holsa")
// 	fmt.Println(val, ok)
// }