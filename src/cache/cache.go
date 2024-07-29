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
	channel := make(chan int)

	for {
		runtime.ReadMemStats(&memStats)

		if memStats.HeapAlloc/1024 > memoryThreshold {
			logging.Warning.Println("Memory usage:", memStats.HeapAlloc/1024, "kB")
			logging.Warning.Println("High memory usage detected, clearing cache")
			go handleHighMemoryUsage(c, channel)
			//Should be blocking until cache is cleared
			entryCleard := <-channel
			logging.Info.Println("Cache cleared, ", entryCleard, " entries removed")
			time.Sleep(2 * time.Second)
			continue
		}
		time.Sleep(1 * time.Second)
	}
}

func handleHighMemoryUsage(cache *Cache, channel chan int) {
	var entriesCleared int
	var minHitCount int
	var maxHitCount int
	var meanHitCount float64
	var nbEntries int
	cache.EntryCache.Range(func(key, value interface{}) bool {
		if minHitCount == 0 || value.(CacheEntry).AccessCount < minHitCount {
			minHitCount = value.(CacheEntry).AccessCount
		}
		if value.(CacheEntry).AccessCount > maxHitCount {
			maxHitCount = value.(CacheEntry).AccessCount
		}
		meanHitCount += float64(value.(CacheEntry).AccessCount)
		nbEntries++
		return true
	})
	meanHitCount = meanHitCount / float64(nbEntries)
	cache.EntryCache.Range(func(key, value interface{}) bool {
		if value.(CacheEntry).AccessCount < int(meanHitCount) {
			cache.EntryCache.Delete(key)
			entriesCleared++
		}
		return true
	})
	channel <- entriesCleared
}
