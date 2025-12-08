package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	data      []byte
}

type Cache struct {
	mu           sync.Mutex
	data         map[string]CacheEntry
	reapInterval time.Duration
}
