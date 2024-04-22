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

