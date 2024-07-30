package cache_test

import (
	"testing"
	"time"

	"github.com/StevenSermeus/goval/src/cache"
)

// TestSetKey tests the SetKey method of the Cache
func TestSetKey(t *testing.T) {
	c := &cache.Cache{}

	key := "testKey"
	value := "testValue"
	valueType := "string"

	c.SetKey(key, value, valueType)

	entry, ok := c.EntryCache.Load(key)
	if !ok {
		t.Fatalf("expected key %s to be present in the cache", key)
	}

	cacheEntry, ok := entry.(cache.CacheEntry)
	if !ok {
		t.Fatalf("expected value to be of type CacheEntry")
	}

	if cacheEntry.Value != value {
		t.Errorf("expected value %s, got %s", value, cacheEntry.Value)
	}

	if cacheEntry.ValueType != valueType {
		t.Errorf("expected value type %s, got %s", valueType, cacheEntry.ValueType)
	}
}

// TestReadKey tests the ReadKey method of the Cache
func TestReadKey(t *testing.T) {
	c := &cache.Cache{}

	key := "testKey"
	value := "testValue"
	valueType := "string"

	c.SetKey(key, value, valueType)

	entry, err := c.ReadKey(key)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if entry.Value != value {
		t.Errorf("expected value %s, got %s", value, entry.Value)
	}

	if entry.ValueType != valueType {
		t.Errorf("expected value type %s, got %s", valueType, entry.ValueType)
	}

	// Verify access count increment and last access time update
	if entry.AccessCount != 1 {
		t.Errorf("expected access count 1, got %d", entry.AccessCount)
	}

	if entry.LastAccess.Before(time.Now().Add(-1 * time.Second)) {
		t.Errorf("expected last access time to be recent")
	}
}

// TestReadKeyNotFound tests the ReadKey method when the key does not exist
func TestReadKeyNotFound(t *testing.T) {
	c := &cache.Cache{}

	key := "nonExistentKey"

	_, err := c.ReadKey(key)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

// TestDeleteKey tests the DeleteKey method of the Cache
func TestDeleteKey(t *testing.T) {
	c := &cache.Cache{}

	key := "testKey"
	value := "testValue"
	valueType := "string"

	c.SetKey(key, value, valueType)
	c.DeleteKey(key)

	_, ok := c.EntryCache.Load(key)
	if ok {
		t.Fatalf("expected key %s to be deleted from the cache", key)
	}
}
