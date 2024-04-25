package cache

import (
	"fmt"
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt 	time.Time
	val			[]byte
}

type Cache struct {
	mu		sync.Mutex
	C 	map[string]CacheEntry
}

func newCacheEntry(val []byte) CacheEntry {
	return CacheEntry{time.Now(), val}
}

func (c Cache) Add(key string, val []byte) {
	fmt.Println("---- Adding cache ----")
	c.mu.Lock()
	defer c.mu.Unlock()
	c.C[key] = newCacheEntry(val)
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, exists := c.C[key]
	return val.val, exists
}

func (c Cache) reapLoop(interval time.Duration) {
	timeNow := time.Now()
	for k, v := range(c.C) {
		if timeNow.Sub(v.createdAt) > interval {
			delete(c.C, k)
		}
	}
}

// func main() {
// 	cache := Cache{cache : make(map[string]CacheEntry)}

// 	val :=  []byte{1,2,3}
// 	cache.Add("hola", val)

// 	val, ok := cache.Get("hola")
// 	fmt.Println(val, ok)

// 	val, ok = cache.Get("holsa")
// 	fmt.Println(val, ok)
// }