package pokecache

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	cache := NewCache(5 * time.Second)
	if cache == nil {
		t.Error("NewCache returned nil")
	}
	if cache.entries == nil {
		t.Error("cache entries map is nil")
	}
	if cache.interval != 5*time.Second {
		t.Errorf("expected interval 5s, got %v", cache.interval)
	}
}

func TestAddAndGet(t *testing.T) {
	cache := NewCache(5 * time.Second)

	key := "test-key"
	val := []byte("test-value")

	cache.Add(key, val)

	retrieved, ok := cache.Get(key)
	if !ok {
		t.Error("expected to find key in cache")
	}

	if string(retrieved) != string(val) {
		t.Errorf("expected %s, got %s", string(val), string(retrieved))
	}
}

func TestGetNonExistent(t *testing.T) {
	cache := NewCache(5 * time.Second)

	_, ok := cache.Get("non-existent-key")
	if ok {
		t.Error("expected false for non-existent key")
	}
}

func TestReap(t *testing.T) {
	interval := 10 * time.Millisecond
	cache := NewCache(interval)

	cache.Add("key1", []byte("value1"))

	// Verify it's there
	_, ok := cache.Get("key1")
	if !ok {
		t.Error("key1 should exist")
	}

	// Wait for reap to happen (interval + small buffer)
	time.Sleep(interval + 15*time.Millisecond)

	// Should be reaped now
	_, ok = cache.Get("key1")
	if ok {
		t.Error("key1 should have been reaped")
	}
}

func TestReapDoesNotRemoveRecent(t *testing.T) {
	interval := 100 * time.Millisecond
	cache := NewCache(interval)

	cache.Add("key1", []byte("value1"))

	// Wait less than the interval
	time.Sleep(50 * time.Millisecond)

	// Should still be there
	_, ok := cache.Get("key1")
	if !ok {
		t.Error("key1 should still exist")
	}
}

func TestMultipleEntries(t *testing.T) {
	cache := NewCache(5 * time.Second)

	cache.Add("key1", []byte("value1"))
	cache.Add("key2", []byte("value2"))
	cache.Add("key3", []byte("value3"))

	val1, ok1 := cache.Get("key1")
	val2, ok2 := cache.Get("key2")
	val3, ok3 := cache.Get("key3")

	if !ok1 || !ok2 || !ok3 {
		t.Error("all keys should exist")
	}

	if string(val1) != "value1" || string(val2) != "value2" || string(val3) != "value3" {
		t.Error("values don't match")
	}
}

func TestOverwriteKey(t *testing.T) {
	cache := NewCache(5 * time.Second)

	cache.Add("key1", []byte("value1"))
	cache.Add("key1", []byte("value2"))

	val, ok := cache.Get("key1")
	if !ok {
		t.Error("key1 should exist")
	}

	if string(val) != "value2" {
		t.Errorf("expected value2, got %s", string(val))
	}
}
