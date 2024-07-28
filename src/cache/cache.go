package cache

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/StevenSermeus/goval/src/logging"
)

// CacheEntry represents a cache entry with a value and last access time.
type CacheEntry struct {
	Value       string
	LastAccess  time.Time
	AccessCount int
}

// Cache is a wrapper around sync.Map for storing CacheEntry values.
type Cache struct {
	EntryCache sync.Map
}

// function to add a key to the cache
func (c *Cache) SetKey(key string, value string) {
	c.EntryCache.Store(key, CacheEntry{Value: value, LastAccess: time.Now(), AccessCount: 0})
}

func (c *Cache) ReadKey(key string) (string, error) {
	value, ok := c.EntryCache.Load(key)
	if !ok {
		fmt.Println("Key not found in cache")
		return "", errors.New("key not found")
	}
	// Update cache entry with last access time
	c.EntryCache.Store(key, CacheEntry{Value: value.(CacheEntry).Value, LastAccess: time.Now(), AccessCount: value.(CacheEntry).AccessCount + 1})
	return value.(CacheEntry).Value, nil
}

func (c *Cache) DeleteKey(key string) {
	c.EntryCache.Delete(key)
}

func (c *Cache) CacheSizeManagement(memoryThreshold uint64) {
	logging.Info.Println("Starting cache size management")
	var memStats runtime.MemStats

	for {
		runtime.ReadMemStats(&memStats)
		if memStats.HeapAlloc/1024/1024 > memoryThreshold {
			logging.Warning.Println("High memory usage detected, clearing cache")
			go handleHighMemoryUsage(c)
			time.Sleep(2 * time.Second)
		}

		time.Sleep(1 * time.Second)
	}
}

func handleHighMemoryUsage(cache *Cache) {
	// For the moment we just clear the cache when the memory threshold is exceeded to free up memory
	cache.EntryCache = sync.Map{}
	fmt.Println("Cache cleared.")
}
